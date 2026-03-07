<template>
  <div>
    <div class="d-flex align-center mb-4">
      <h1 class="text-h4">{{ isEdit ? 'Edit Document' : 'New Document' }}</h1>
      <v-spacer />
      <v-menu v-if="isSaved">
        <template v-slot:activator="{ props }">
          <v-btn icon variant="text" v-bind="props">
            <v-icon>mdi-dots-vertical</v-icon>
          </v-btn>
        </template>
        <v-list>
          <v-list-item @click="showDeleteConfirm = true">
            <template v-slot:prepend>
              <v-icon color="error">mdi-delete</v-icon>
            </template>
            <v-list-item-title class="text-error">Send to Trash</v-list-item-title>
          </v-list-item>
        </v-list>
      </v-menu>
      <v-btn variant="text" @click="goBack">{{ isSaved ? 'Close' : 'Cancel' }}</v-btn>
      <v-btn color="primary" :loading="saving" class="ml-2" @click="save">Save</v-btn>

    <v-dialog v-model="showDeleteConfirm" max-width="400">
      <v-card>
        <v-card-title>Send to Trash</v-card-title>
        <v-card-text>Are you sure you want to send this document to the trash?</v-card-text>
        <v-card-actions>
          <v-spacer />
          <v-btn @click="showDeleteConfirm = false">Cancel</v-btn>
          <v-btn color="error" @click="deleteDoc">Send to Trash</v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>
    </div>

    <v-progress-linear v-if="loading" indeterminate color="primary" />

    <template v-if="!loading">
      <v-row class="mb-3">
        <v-col cols="12" md="6">
          <v-text-field v-model="form.title" label="Title" required hide-details />
        </v-col>
        <v-col cols="12" md="3">
          <v-select
            v-model="form.category_id"
            label="Category"
            :items="categoryItems"
            item-title="name"
            item-value="id"
            hide-details
          />
        </v-col>
        <v-col cols="12" md="3">
          <v-select
            v-model="form.status"
            label="Status"
            :items="['draft', 'complete']"
            hide-details
          />
        </v-col>
      </v-row>

      <v-card variant="outlined" class="editor-card">
        <div v-if="editor" class="editor-toolbar pa-1 d-flex flex-wrap ga-1">
          <v-btn-group density="compact" variant="text">
            <v-btn size="small" :class="{ 'text-primary': editor.isActive('bold') }" @click="editor.chain().focus().toggleBold().run()">
              <v-icon>mdi-format-bold</v-icon>
            </v-btn>
            <v-btn size="small" :class="{ 'text-primary': editor.isActive('italic') }" @click="editor.chain().focus().toggleItalic().run()">
              <v-icon>mdi-format-italic</v-icon>
            </v-btn>
            <v-btn size="small" :class="{ 'text-primary': editor.isActive('underline') }" @click="editor.chain().focus().toggleUnderline().run()">
              <v-icon>mdi-format-underline</v-icon>
            </v-btn>
            <v-btn size="small" :class="{ 'text-primary': editor.isActive('strike') }" @click="editor.chain().focus().toggleStrike().run()">
              <v-icon>mdi-format-strikethrough</v-icon>
            </v-btn>
          </v-btn-group>

          <v-divider vertical class="mx-1" />

          <v-btn-group density="compact" variant="text">
            <v-btn size="small" :class="{ 'text-primary': editor.isActive('heading', { level: 1 }) }" @click="editor.chain().focus().toggleHeading({ level: 1 }).run()">
              H1
            </v-btn>
            <v-btn size="small" :class="{ 'text-primary': editor.isActive('heading', { level: 2 }) }" @click="editor.chain().focus().toggleHeading({ level: 2 }).run()">
              H2
            </v-btn>
            <v-btn size="small" :class="{ 'text-primary': editor.isActive('heading', { level: 3 }) }" @click="editor.chain().focus().toggleHeading({ level: 3 }).run()">
              H3
            </v-btn>
          </v-btn-group>

          <v-divider vertical class="mx-1" />

          <v-btn-group density="compact" variant="text">
            <v-btn size="small" :class="{ 'text-primary': editor.isActive('bulletList') }" @click="editor.chain().focus().toggleBulletList().run()">
              <v-icon>mdi-format-list-bulleted</v-icon>
            </v-btn>
            <v-btn size="small" :class="{ 'text-primary': editor.isActive('orderedList') }" @click="editor.chain().focus().toggleOrderedList().run()">
              <v-icon>mdi-format-list-numbered</v-icon>
            </v-btn>
            <v-btn size="small" :class="{ 'text-primary': editor.isActive('blockquote') }" @click="editor.chain().focus().toggleBlockquote().run()">
              <v-icon>mdi-format-quote-close</v-icon>
            </v-btn>
          </v-btn-group>

          <v-divider vertical class="mx-1" />

          <v-btn-group density="compact" variant="text">
            <v-btn size="small" @click="editor.chain().focus().setHorizontalRule().run()">
              <v-icon>mdi-minus</v-icon>
            </v-btn>
            <v-btn size="small" :class="{ 'text-primary': editor.isActive('code') }" @click="editor.chain().focus().toggleCode().run()">
              <v-icon>mdi-code-tags</v-icon>
            </v-btn>
            <v-btn size="small" :class="{ 'text-primary': editor.isActive('codeBlock') }" @click="editor.chain().focus().toggleCodeBlock().run()">
              <v-icon>mdi-code-braces</v-icon>
            </v-btn>
          </v-btn-group>

          <v-divider vertical class="mx-1" />

          <v-btn-group density="compact" variant="text">
            <v-btn size="small" @click="editor.chain().focus().undo().run()" :disabled="!editor.can().undo()">
              <v-icon>mdi-undo</v-icon>
            </v-btn>
            <v-btn size="small" @click="editor.chain().focus().redo().run()" :disabled="!editor.can().redo()">
              <v-icon>mdi-redo</v-icon>
            </v-btn>
          </v-btn-group>
        </div>

        <v-divider />

        <editor-content :editor="editor" class="editor-content" />
      </v-card>

      <v-textarea
        v-model="form.secure_notes"
        label="Confidential Notes"
        rows="2"
        prepend-inner-icon="mdi-lock"
        hint="Will NOT appear on the printed cover letter."
        persistent-hint
        class="mt-4"
      />
    </template>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted, onBeforeUnmount, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useDocumentStore } from '../stores/documents'
import { useDocumentCategoryStore } from '../stores/documentCategories'
import { useEditor, EditorContent } from '@tiptap/vue-3'
import StarterKit from '@tiptap/starter-kit'
import Underline from '@tiptap/extension-underline'
import Link from '@tiptap/extension-link'
import { Markdown } from 'tiptap-markdown'

const route = useRoute()
const router = useRouter()
const documentStore = useDocumentStore()
const categoryStore = useDocumentCategoryStore()

const loading = ref(false)
const saving = ref(false)
const showDeleteConfirm = ref(false)

const isEdit = computed(() => !!route.params.id)
const isSaved = computed(() => !!route.params.id)

const categoryItems = computed(() => categoryStore.categories)

const form = reactive({
  title: '',
  category_id: null as number | null,
  content: '',
  status: 'draft' as 'draft' | 'complete',
  secure_notes: '',
})

const editor = useEditor({
  extensions: [
    StarterKit,
    Underline,
    Link.configure({ openOnClick: false }),
    Markdown,
  ],
  content: '',
  onUpdate: ({ editor: e }) => {
    form.content = e.storage.markdown.getMarkdown()
  },
})

watch(() => form.content, (val) => {
  if (editor.value && editor.value.storage.markdown.getMarkdown() !== val) {
    editor.value.commands.setContent(val)
  }
})

function goBack() {
  router.push({ name: 'documents' })
}

async function save() {
  saving.value = true
  try {
    const payload: Record<string, any> = {
      ...form,
      doc_type: 'typed',
    }
    if (isEdit.value) {
      await documentStore.updateDocument(Number(route.params.id), payload)
    } else {
      const created = await documentStore.createDocument(payload)
      router.replace({ name: 'document-edit', params: { id: created.id } })
    }
  } finally {
    saving.value = false
  }
}

async function deleteDoc() {
  showDeleteConfirm.value = false
  await documentStore.deleteDocument(Number(route.params.id))
  router.push({ name: 'documents' })
}

onMounted(async () => {
  categoryStore.fetchCategories()

  if (isEdit.value) {
    loading.value = true
    try {
      const doc = await documentStore.fetchDocument(Number(route.params.id))
      form.title = doc.title
      form.category_id = doc.category_id
      form.content = doc.content || ''
      form.status = doc.status
      form.secure_notes = doc.secure_notes || ''
      if (editor.value) {
        editor.value.commands.setContent(form.content)
      }
    } finally {
      loading.value = false
    }
  }
})

onBeforeUnmount(() => {
  editor.value?.destroy()
})
</script>

<style scoped>
.editor-card {
  min-height: 50vh;
  display: flex;
  flex-direction: column;
}

.editor-toolbar {
  background: rgb(var(--v-theme-surface));
  border-bottom: none;
}

.editor-content {
  flex: 1;
  padding: 16px;
  min-height: 45vh;
}

.editor-content :deep(.tiptap) {
  outline: none;
  min-height: 40vh;
}

.editor-content :deep(.tiptap p) {
  margin-bottom: 0.75em;
}

.editor-content :deep(.tiptap h1),
.editor-content :deep(.tiptap h2),
.editor-content :deep(.tiptap h3) {
  margin-top: 1em;
  margin-bottom: 0.5em;
}

.editor-content :deep(.tiptap ul),
.editor-content :deep(.tiptap ol) {
  padding-left: 1.5em;
  margin-bottom: 0.75em;
}

.editor-content :deep(.tiptap blockquote) {
  border-left: 3px solid rgb(var(--v-theme-primary));
  padding-left: 1em;
  margin-left: 0;
  margin-bottom: 0.75em;
  color: rgba(var(--v-theme-on-surface), 0.7);
}

.editor-content :deep(.tiptap pre) {
  background: rgba(var(--v-theme-on-surface), 0.05);
  border-radius: 4px;
  padding: 0.75em 1em;
  margin-bottom: 0.75em;
}

.editor-content :deep(.tiptap code) {
  background: rgba(var(--v-theme-on-surface), 0.08);
  border-radius: 3px;
  padding: 0.15em 0.3em;
  font-size: 0.9em;
}

.editor-content :deep(.tiptap hr) {
  border: none;
  border-top: 1px solid rgba(var(--v-theme-on-surface), 0.12);
  margin: 1em 0;
}
</style>
