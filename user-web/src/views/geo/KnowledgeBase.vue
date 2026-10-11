<template>
  <div class="geo-page">

    <div class="p-4" style="display:flex; gap:16px; min-height:calc(100vh - 180px)">
    
    <div style="flex:1; min-width:0; display:flex; flex-direction:column; gap:16px">
      <el-row :gutter="16">
        <el-col :span="8">
          <el-card><el-statistic :title="$t('auto.KnowledgeBase.1.d095d1')" :value="docs.length" /></el-card>
        </el-col>
        <el-col :span="8">
          <el-card><el-statistic :title="$t('auto.KnowledgeBase.2.8e7d4a')" :value="aLevelCount" /></el-card>
        </el-col>
        <el-col :span="8">
          <el-card><el-statistic :title="$t('auto.KnowledgeBase.3.05f8c3')" :value="totalEntities" /></el-card>
        </el-col>
      </el-row>

      <el-card style="flex:1">
        <template #header>
          <div class="flex items-center justify-between">
            <span class="font-bold">{{ $t('auto.KnowledgeBase.11.a4005a') }}</span>
            <div class="flex items-center gap-2">
              <el-input v-model="keyword" :placeholder="$t('auto.KnowledgeBase.4.cebdbd')" size="small" style="width:220px" clearable @clear="loadDocs" @keyup.enter="loadDocs" />
              <el-upload :show-file-list="false" :before-upload="onUpload" :http-request="() => {}">
                <el-button size="small" type="primary">{{ $t('auto.KnowledgeBase.12.c485b3') }}</el-button>
              </el-upload>
            </div>
          </div>
        </template>
        <el-table :data="docs" v-loading="loading" size="small">
          <el-table-column prop="title" :label="$t('auto.KnowledgeBase.5.325369')" min-width="200">
            <template #default="{ row }">
              <div class="font-medium">{{ row.title || row.filename }}</div>
              <div class="text-xs text-gray-500">{{ row.doc_id || row.id }}</div>
            </template>
          </el-table-column>
          <el-table-column :label="$t('auto.KnowledgeBase.6.95e0d7')" width="140">
            <template #default="{ row }">
              <el-select v-model="row.source_level" size="small" style="width:90px" @change="onUpdateLevel(row)">
                <el-option label="A" value="A" />
                <el-option label="B" value="B" />
                <el-option label="C" value="C" />
                <el-option label="D" value="D" />
              </el-select>
            </template>
          </el-table-column>
          <el-table-column prop="entity_count" :label="$t('auto.KnowledgeBase.7.cc8f47')" width="90" />
          <el-table-column prop="uploaded_at" :label="$t('auto.KnowledgeBase.8.cae255')" width="170" />
          <el-table-column :label="$t('auto.KnowledgeBase.9.2b6bc0')" width="200">
            <template #default="{ row }">
              <el-button size="small" type="primary" plain @click="onExtract(row)" :loading="extractLoading === row.id">{{ $t('auto.KnowledgeBase.13.cdc0cc') }}</el-button>
              <el-button size="small" type="danger" text @click="onDelete(row)">{{ $t('auto.KnowledgeBase.14.2f4aad') }}</el-button>
            </template>
          </el-table-column>
        </el-table>
      </el-card>
    </div>

    
    <el-card style="width:400px; flex-shrink:0; display:flex; flex-direction:column">
      <template #header><span class="font-bold">Ask GEO KB</span></template>
      <div class="chat-panel" style="flex:1; overflow-y:auto; padding-right:4px">
        <div v-for="(m, i) in messages" :key="i" class="chat-message" :class="m.role">
          <div class="chat-avatar">{{ m.role === 'user' ? '我' : 'AI' }}</div>
          <div class="chat-bubble">{{ m.content }}</div>
        </div>
        <div v-if="messages.length === 0" class="text-gray-400 text-sm text-center py-4">{{ $t('auto.KnowledgeBase.15.085846') }}</div>
      </div>
      <div class="mt-3">
        <el-input v-model="askInput" type="textarea" :rows="3" :placeholder="$t('auto.KnowledgeBase.10.c32cbb')" @keyup.enter.ctrl="onAsk" resize="none" />
        <div class="flex justify-end mt-2">
          <el-button type="primary" :loading="asking" :disabled="!askInput.trim()" @click="onAsk">{{ $t('auto.KnowledgeBase.16.1535fc') }}</el-button>
        </div>
      </div>
    </el-card>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { ElMessage } from 'element-plus'
import { listKBDocuments, saveKBDocument, deleteKBDocument, askKB, extractEntities } from '@/api/geoEntity.js'

const { t } = useI18n()

const docs = ref([])
const loading = ref(false)
const keyword = ref('')
const extractLoading = ref(null)

const aLevelCount = computed(() => docs.value.filter(d => d.source_level === 'A').length)
const totalEntities = computed(() => docs.value.reduce((s, d) => s + (d.entity_count || 0), 0))

const messages = ref([])
const askInput = ref('')
const asking = ref(false)

const loadDocs = async () => {
  loading.value = true
  try {
    const data = await listKBDocuments(keyword.value)
    docs.value = Array.isArray(data) ? data : (data?.list || data?.items || [])
  } catch { docs.value = [] }
  loading.value = false
}

const onUpload = async (file) => {
  const rawFile = file.raw || file
  const ext = (rawFile.name || '').split('.').pop()?.toLowerCase() || '';
  const isText = ['txt', 'md', 'markdown', 'json'].includes(ext) ||
    rawFile.type?.startsWith('text/') ||
    rawFile.type === 'application/json'
  if (!isText) {
    ElMessage.warning(t('auto.KnowledgeBase.17.4d5b8a'))
    return false
  }
  try {
    const content = await new Promise((resolve, reject) => {
      const reader = new FileReader()
      reader.onload = () => resolve(String(reader.result || ''))
      reader.onerror = () => reject(new Error(t('auto.KnowledgeBase.18.e2bb17')))
      reader.readAsText(rawFile)
    })
    if (!content.trim()) {
      ElMessage.warning(t('auto.KnowledgeBase.19.aa6673'))
      return false
    }
    const title = rawFile.name.replace(/\.[^.]+$/, '') || t('auto.KnowledgeBase.20.1232aa')
    await saveKBDocument({ title, content, source_level: 'C' })
    ElMessage.success(t('auto.KnowledgeBase.21.da1425'))
    await loadDocs()
  } catch (e) {
    ElMessage.error(e?.message || t('auto.KnowledgeBase.22.54e5de'))
  }
  return false
}

const onUpdateLevel = async (row) => {
  try {
    await saveKBDocument({ id: row.id || row.doc_id, source_level: row.source_level, title: row.title || row.filename, content: row.content || '' })
    ElMessage.success(t('auto.KnowledgeBase.23.1b9209'))
  } catch (e) { ElMessage.error(e?.message || t('auto.KnowledgeBase.24.930442')) }
}

const onDelete = async (row) => {
  try {
    await deleteKBDocument(row.id || row.doc_id)
    ElMessage.success(t('auto.KnowledgeBase.25.5cc232'))
    loadDocs()
  } catch (e) { ElMessage.error(e?.message || t('auto.KnowledgeBase.26.acf066')) }
}

const onExtract = async (row) => {
  extractLoading.value = row.id
  try {
    await extractEntities(row.id || row.doc_id)
    ElMessage.success(t('auto.KnowledgeBase.27.af2adb'))
    loadDocs()
  } catch (e) { ElMessage.error(t('auto.KnowledgeBase.28.38bca1') + (e?.message || e)) }
  extractLoading.value = null
}

const onAsk = async () => {
  const q = askInput.value.trim()
  if (!q) return
  messages.value.push({ role: 'user', content: q })
  askInput.value = ''
  asking.value = true
  try {
    const answer = await askKB(q)
    const content = typeof answer === 'string' ? answer : (answer?.answer || answer?.data || JSON.stringify(answer))
    messages.value.push({ role: 'assistant', content })
  } catch (e) {
    messages.value.push({ role: 'assistant', content: t('auto.KnowledgeBase.29.5a2b64') + (e?.message || e) })
  }
  asking.value = false
}

onMounted(loadDocs)
</script>

<style lang="scss" scoped>
.chat-panel {
  display: flex;
  flex-direction: column;
  gap: $spacing-md;
  min-height: 300px;
}
.chat-message {
  display: flex;
  gap: $spacing-sm;
  align-items: flex-start;
}
.chat-message.user { flex-direction: row-reverse; }
.chat-avatar {
  width: 28px; height: 28px; border-radius: 50%;
  background: var(--el-color-primary); color: $bg-color;
  display: flex; align-items: center; justify-content: center;
  font-size: $font-size-extra-small; flex-shrink: 0;
}
.chat-message.assistant .chat-avatar { background: $success-color; }
.chat-bubble {
  max-width: 75%;
  padding: $spacing-md 14px;
  border-radius: 10px;
  background: var(--el-fill-color-light);
  font-size: $font-size-base;
  white-space: pre-wrap;
  word-break: break-word;
}
.chat-message.user .chat-bubble {
  background: var(--el-color-primary-light-9);
  color: var(--el-text-color-primary);
}
</style>
