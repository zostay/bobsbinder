package services

import (
	"database/sql"
	"fmt"
	"strings"

	"go.uber.org/zap"

	"github.com/zostay/bobsbinder/internal/models"
)

type sectionDef struct {
	Key   string
	Title string
}

var defaultSections = []sectionDef{
	{Key: "contacts", Title: "Contacts"},
	{Key: "documents", Title: "Documents"},
	{Key: "locations", Title: "Important Locations"},
	{Key: "digital_info", Title: "Digital Information"},
}

type computedItem struct {
	SourceType string
	SourceID   int64
	Content    string
	ItemType   string
}

// SyncLetter ensures the survivor letter exists for the given user, syncs all
// sections and auto-generated items from source data, and returns the full letter.
func SyncLetter(db *sql.DB, logger *zap.Logger, userID int64) (*models.FullSurvivorLetter, error) {
	// 1. Ensure letter row exists (with personalized defaults from user name)
	var userName string
	if err := db.QueryRow("SELECT name FROM users WHERE id = ?", userID).Scan(&userName); err != nil {
		userName = ""
	}

	defaultGreeting := "Dear loved ones,"
	defaultIntro := "If you are reading this, I am no longer able to manage my own affairs. I have prepared this letter to help you find and manage the important things in my life. Please use this as a guide — it covers the key documents, accounts, contacts, and other details you may need."
	defaultClosing := "I hope this letter makes a difficult time a little easier. Please know that I love you and I am grateful for everything you have meant to me."
	defaultSignature := userName

	_, err := db.Exec(
		"INSERT IGNORE INTO survivor_letters (user_id, greeting, intro, closing, signature) VALUES (?, ?, ?, ?, ?)",
		userID, defaultGreeting, defaultIntro, defaultClosing, defaultSignature,
	)
	if err != nil {
		return nil, fmt.Errorf("ensure letter: %w", err)
	}

	var letter models.SurvivorLetter
	err = db.QueryRow(`
		SELECT id, user_id, greeting, intro, closing, signature, created_at, updated_at
		FROM survivor_letters WHERE user_id = ?
	`, userID).Scan(&letter.ID, &letter.UserID, &letter.Greeting, &letter.Intro,
		&letter.Closing, &letter.Signature, &letter.CreatedAt, &letter.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("fetch letter: %w", err)
	}

	// 2. Ensure all standard sections exist
	for i, sec := range defaultSections {
		_, err := db.Exec(`
			INSERT INTO survivor_letter_sections (letter_id, section_key, title, sort_order)
			VALUES (?, ?, ?, ?)
			ON DUPLICATE KEY UPDATE letter_id = letter_id
		`, letter.ID, sec.Key, sec.Title, i+1)
		if err != nil {
			return nil, fmt.Errorf("ensure section %s: %w", sec.Key, err)
		}
	}

	// 3. Load all sections
	sectionRows, err := db.Query(`
		SELECT id, letter_id, section_key, title, sort_order, visible, created_at, updated_at
		FROM survivor_letter_sections WHERE letter_id = ? ORDER BY sort_order
	`, letter.ID)
	if err != nil {
		return nil, fmt.Errorf("load sections: %w", err)
	}
	defer sectionRows.Close()

	var sections []models.SurvivorLetterSection
	sectionMap := make(map[string]int64) // section_key -> section_id
	for sectionRows.Next() {
		var s models.SurvivorLetterSection
		if err := sectionRows.Scan(&s.ID, &s.LetterID, &s.SectionKey, &s.Title,
			&s.SortOrder, &s.Visible, &s.CreatedAt, &s.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scan section: %w", err)
		}
		s.Items = []models.SurvivorLetterItem{}
		sections = append(sections, s)
		sectionMap[s.SectionKey] = s.ID
	}

	// 4. Sync items for each section
	for sectionKey, sectionID := range sectionMap {
		items, err := computeItems(db, userID, sectionKey)
		if err != nil {
			logger.Error("failed to compute items for section", zap.String("section", sectionKey), zap.Error(err))
			continue
		}

		if err := syncSectionItems(db, sectionID, items); err != nil {
			logger.Error("failed to sync items for section", zap.String("section", sectionKey), zap.Error(err))
		}
	}

	// 5. Load all items for all sections
	for i := range sections {
		itemRows, err := db.Query(`
			SELECT id, section_id, source_type, source_id, content, item_type, provenance, suppressed, sort_order, created_at, updated_at
			FROM survivor_letter_items WHERE section_id = ? ORDER BY sort_order, id
		`, sections[i].ID)
		if err != nil {
			return nil, fmt.Errorf("load items for section %s: %w", sections[i].SectionKey, err)
		}

		for itemRows.Next() {
			var item models.SurvivorLetterItem
			var sourceType sql.NullString
			var sourceID sql.NullInt64
			if err := itemRows.Scan(&item.ID, &item.SectionID, &sourceType, &sourceID,
				&item.Content, &item.ItemType, &item.Provenance, &item.Suppressed,
				&item.SortOrder, &item.CreatedAt, &item.UpdatedAt); err != nil {
				itemRows.Close()
				return nil, fmt.Errorf("scan item: %w", err)
			}
			if sourceType.Valid {
				item.SourceType = &sourceType.String
			}
			if sourceID.Valid {
				item.SourceID = &sourceID.Int64
			}
			sections[i].Items = append(sections[i].Items, item)
		}
		itemRows.Close()
	}

	return &models.FullSurvivorLetter{
		SurvivorLetter: letter,
		Sections:       sections,
	}, nil
}

func computeItems(db *sql.DB, userID int64, sectionKey string) ([]computedItem, error) {
	switch sectionKey {
	case "contacts":
		return computeContacts(db, userID)
	case "documents":
		return computeDocumentsAll(db, userID)
	case "locations":
		return computeLocationsAll(db, userID)
	case "digital_info":
		return computeDigitalInfoAll(db, userID)
	default:
		return nil, nil
	}
}

func computeDocumentsAll(db *sql.DB, userID int64) ([]computedItem, error) {
	allSlugs := []string{
		"will", "final-arrangements", "other",
		"poa-medical", "poa-financial", "medical-directives",
		"pre-certification", "memorial", "remains",
	}
	docs, err := computeDocumentsByCategory(db, userID, allSlugs)
	if err != nil {
		return nil, err
	}

	policies, err := computeInsurancePolicies(db, userID)
	if err != nil {
		return nil, err
	}

	obituary, err := computeObituaryInfo(db, userID)
	if err != nil {
		return nil, err
	}

	items := append(docs, policies...)
	items = append(items, obituary...)
	return items, nil
}

func computeLocationsAll(db *sql.DB, userID int64) ([]computedItem, error) {
	physical, err := computeLocations(db, userID, "physical")
	if err != nil {
		return nil, err
	}

	digital, err := computeLocations(db, userID, "digital")
	if err != nil {
		return nil, err
	}

	return append(physical, digital...), nil
}

func computeDigitalInfoAll(db *sql.DB, userID int64) ([]computedItem, error) {
	access, err := computeDigitalAccess(db, userID, []string{"computer", "phone", "password_manager"})
	if err != nil {
		return nil, err
	}

	financial, err := computeServiceAccounts(db, userID, "financial_tool")
	if err != nil {
		return nil, err
	}

	backup, err := computeServiceAccounts(db, userID, "backup_service")
	if err != nil {
		return nil, err
	}

	taxPrep, err := computeServiceAccounts(db, userID, "tax_preparer")
	if err != nil {
		return nil, err
	}

	items := append(access, financial...)
	items = append(items, backup...)
	items = append(items, taxPrep...)
	return items, nil
}

func computeContacts(db *sql.DB, userID int64) ([]computedItem, error) {
	rows, err := db.Query(`SELECT id, name, role, phone, email FROM contacts WHERE user_id = ? ORDER BY is_primary DESC, name`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []computedItem
	for rows.Next() {
		var id int64
		var name, role, phone, email string
		if err := rows.Scan(&id, &name, &role, &phone, &email); err != nil {
			continue
		}
		parts := []string{name}
		if role != "" {
			parts = append(parts, "("+role+")")
		}
		if phone != "" {
			parts = append(parts, "- "+phone)
		}
		if email != "" {
			parts = append(parts, "- "+email)
		}
		items = append(items, computedItem{
			SourceType: "contact",
			SourceID:   id,
			Content:    strings.Join(parts, " "),
			ItemType:   "numbered",
		})
	}
	return items, nil
}

func computeDocumentsByCategory(db *sql.DB, userID int64, categorySlugs []string) ([]computedItem, error) {
	if len(categorySlugs) == 0 {
		return nil, nil
	}
	placeholders := make([]string, len(categorySlugs))
	args := make([]any, 0, len(categorySlugs)+1)
	args = append(args, userID)
	for i, slug := range categorySlugs {
		placeholders[i] = "?"
		args = append(args, slug)
	}

	query := fmt.Sprintf(`
		SELECT d.id, d.title, dc.name as category_name, l.name as location_name
		FROM documents d
		JOIN parties p ON d.party_id = p.id
		JOIN document_categories dc ON d.category_id = dc.id
		LEFT JOIN locations l ON d.location_id = l.id
		WHERE p.user_id = ? AND dc.slug IN (%s)
		ORDER BY d.title
	`, strings.Join(placeholders, ","))

	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []computedItem
	for rows.Next() {
		var id int64
		var title, categoryName string
		var locationName sql.NullString
		if err := rows.Scan(&id, &title, &categoryName, &locationName); err != nil {
			continue
		}
		content := fmt.Sprintf("%s (%s)", title, categoryName)
		if locationName.Valid && locationName.String != "" {
			content += fmt.Sprintf(" - located at %s", locationName.String)
		}
		items = append(items, computedItem{
			SourceType: "document",
			SourceID:   id,
			Content:    content,
			ItemType:   "numbered",
		})
	}
	return items, nil
}

func computeDigitalAccess(db *sql.DB, userID int64, types []string) ([]computedItem, error) {
	if len(types) == 0 {
		return nil, nil
	}
	placeholders := make([]string, len(types))
	args := make([]any, 0, len(types)+1)
	args = append(args, userID)
	for i, t := range types {
		placeholders[i] = "?"
		args = append(args, t)
	}

	query := fmt.Sprintf(`
		SELECT da.id, da.name, da.type, da.username, da.instructions
		FROM digital_access da
		WHERE da.user_id = ? AND da.type IN (%s)
		ORDER BY da.name
	`, strings.Join(placeholders, ","))

	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []computedItem
	for rows.Next() {
		var id int64
		var name, daType, username, instructions string
		if err := rows.Scan(&id, &name, &daType, &username, &instructions); err != nil {
			continue
		}
		content := name
		if username != "" {
			content += fmt.Sprintf(" (user: %s)", username)
		}
		if instructions != "" {
			content += " - " + instructions
		}
		items = append(items, computedItem{
			SourceType: "digital_access",
			SourceID:   id,
			Content:    content,
			ItemType:   "numbered",
		})
	}
	return items, nil
}

func computeInsurancePolicies(db *sql.DB, userID int64) ([]computedItem, error) {
	rows, err := db.Query(`
		SELECT ip.id, ip.provider, ip.policy_number, ip.type, ip.coverage_amount, ip.beneficiary, ip.agent_name, ip.agent_phone
		FROM insurance_policies ip
		WHERE ip.user_id = ?
		ORDER BY ip.provider
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []computedItem
	for rows.Next() {
		var id int64
		var provider, policyNumber, pType, beneficiary, agentName, agentPhone string
		var coverageAmount sql.NullFloat64
		if err := rows.Scan(&id, &provider, &policyNumber, &pType, &coverageAmount, &beneficiary, &agentName, &agentPhone); err != nil {
			continue
		}
		content := provider
		if pType != "" {
			content += fmt.Sprintf(" (%s)", pType)
		}
		if policyNumber != "" {
			content += fmt.Sprintf(" #%s", policyNumber)
		}
		if coverageAmount.Valid {
			content += fmt.Sprintf(" - $%.2f", coverageAmount.Float64)
		}
		if beneficiary != "" {
			content += fmt.Sprintf(", beneficiary: %s", beneficiary)
		}
		if agentName != "" {
			content += fmt.Sprintf(", agent: %s", agentName)
			if agentPhone != "" {
				content += fmt.Sprintf(" (%s)", agentPhone)
			}
		}
		items = append(items, computedItem{
			SourceType: "insurance_policy",
			SourceID:   id,
			Content:    content,
			ItemType:   "bulleted",
		})
	}
	return items, nil
}

func computeLocations(db *sql.DB, userID int64, locType string) ([]computedItem, error) {
	rows, err := db.Query(`
		SELECT id, name, description, address, access_instructions
		FROM locations WHERE user_id = ? AND type = ? ORDER BY name
	`, userID, locType)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []computedItem
	for rows.Next() {
		var id int64
		var name, description, address, accessInstructions string
		if err := rows.Scan(&id, &name, &description, &address, &accessInstructions); err != nil {
			continue
		}
		content := name
		if description != "" {
			content += " - " + description
		}
		if address != "" {
			content += fmt.Sprintf(" (%s)", address)
		}
		if accessInstructions != "" {
			content += " [" + accessInstructions + "]"
		}
		items = append(items, computedItem{
			SourceType: "location",
			SourceID:   id,
			Content:    content,
			ItemType:   "numbered",
		})
	}
	return items, nil
}

func computeServiceAccounts(db *sql.DB, userID int64, saType string) ([]computedItem, error) {
	rows, err := db.Query(`
		SELECT id, name, provider, account_number, contact_name, contact_phone, contact_email, notes
		FROM service_accounts WHERE user_id = ? AND type = ? ORDER BY name
	`, userID, saType)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []computedItem
	for rows.Next() {
		var id int64
		var name, provider, accountNumber, contactName, contactPhone, contactEmail, notes string
		if err := rows.Scan(&id, &name, &provider, &accountNumber, &contactName, &contactPhone, &contactEmail, &notes); err != nil {
			continue
		}
		content := name
		if provider != "" {
			content += " (" + provider + ")"
		}
		if accountNumber != "" {
			content += " #" + accountNumber
		}
		if contactName != "" {
			content += fmt.Sprintf(", contact: %s", contactName)
			if contactPhone != "" {
				content += " " + contactPhone
			}
		}
		items = append(items, computedItem{
			SourceType: "service_account",
			SourceID:   id,
			Content:    content,
			ItemType:   "numbered",
		})
	}
	return items, nil
}

func computeObituaryInfo(db *sql.DB, userID int64) ([]computedItem, error) {
	rows, err := db.Query(`
		SELECT poi.id, poi.type, poi.name, poi.relationship, poi.details, p.name as party_name
		FROM party_obituary_info poi
		JOIN parties p ON poi.party_id = p.id
		WHERE p.user_id = ?
		ORDER BY poi.type, poi.name
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []computedItem
	for rows.Next() {
		var id int64
		var oType, name, relationship, details, partyName string
		if err := rows.Scan(&id, &oType, &name, &relationship, &details, &partyName); err != nil {
			continue
		}
		content := name
		if relationship != "" {
			content += fmt.Sprintf(" (%s)", relationship)
		}
		if details != "" {
			content += " - " + details
		}
		items = append(items, computedItem{
			SourceType: "obituary_info",
			SourceID:   id,
			Content:    content,
			ItemType:   "bulleted",
		})
	}
	return items, nil
}

func syncSectionItems(db *sql.DB, sectionID int64, computed []computedItem) error {
	// Build a map of existing auto items by (source_type, source_id)
	existingRows, err := db.Query(`
		SELECT id, source_type, source_id, content, provenance, suppressed
		FROM survivor_letter_items
		WHERE section_id = ? AND source_type IS NOT NULL
	`, sectionID)
	if err != nil {
		return err
	}
	defer existingRows.Close()

	type existingItem struct {
		ID         int64
		Content    string
		Provenance string
		Suppressed bool
	}
	existing := make(map[string]existingItem) // "source_type:source_id" -> item
	for existingRows.Next() {
		var id, sourceID int64
		var sourceType, content, provenance string
		var suppressed bool
		if err := existingRows.Scan(&id, &sourceType, &sourceID, &content, &provenance, &suppressed); err != nil {
			continue
		}
		key := fmt.Sprintf("%s:%d", sourceType, sourceID)
		existing[key] = existingItem{ID: id, Content: content, Provenance: provenance, Suppressed: suppressed}
	}

	// Track which existing items are still present in computed
	seen := make(map[string]bool)

	// Get max sort_order for new items
	var maxSort int
	db.QueryRow("SELECT COALESCE(MAX(sort_order), 0) FROM survivor_letter_items WHERE section_id = ?", sectionID).Scan(&maxSort)

	for _, ci := range computed {
		key := fmt.Sprintf("%s:%d", ci.SourceType, ci.SourceID)
		seen[key] = true

		if ex, ok := existing[key]; ok {
			// Item exists
			if ex.Suppressed {
				continue // user suppressed it, skip
			}
			if ex.Provenance == "auto" && ex.Content != ci.Content {
				// Update auto item with new content
				db.Exec(`UPDATE survivor_letter_items SET content = ?, item_type = ? WHERE id = ?`,
					ci.Content, ci.ItemType, ex.ID)
			}
			// auto_edited or manual: skip
		} else {
			// New item
			maxSort++
			db.Exec(`
				INSERT INTO survivor_letter_items (section_id, source_type, source_id, content, item_type, provenance, sort_order)
				VALUES (?, ?, ?, ?, ?, 'auto', ?)
			`, sectionID, ci.SourceType, ci.SourceID, ci.Content, ci.ItemType, maxSort)
		}
	}

	// Clean up orphans: auto items whose source no longer exists
	for key, ex := range existing {
		if !seen[key] {
			if ex.Provenance == "auto" {
				db.Exec("DELETE FROM survivor_letter_items WHERE id = ?", ex.ID)
			} else if ex.Provenance == "auto_edited" {
				// Convert to manual so user edit is preserved
				db.Exec("UPDATE survivor_letter_items SET provenance = 'manual', source_type = NULL, source_id = NULL WHERE id = ?", ex.ID)
			}
		}
	}

	return nil
}
