<template>
  <div class="batch-import-page">
    <el-alert
      ::title="$t('auto.BatchImport.af030d46')"
      type="info"
      :closable="false"
      show-icon
    >
      <template #default>
        支持两种方式:
        <strong>{{ $t('auto.BatchImport.ca33085f') }}</strong> CSV / JSON / Excel(纯文本),
        <strong>{{ $t('auto.BatchImport.29b14ae0') }}</strong> JSON 数组直接提交。
        每行记录包含 title、content、category、tags 四个字段(除 content 外均可选)。
      </template>
    </el-alert>

    <el-row :gutter="16" style="margin-top: 16px">
      
      <el-col :span="14">
        <el-card>
          <el-tabs v-model="activeTab" @tab-change="handleTabChange">
            
            <el-tab-pane :label="$t('auto.BatchImport.481e034e')" name="upload">
              <el-form label-width="100px" :model="uploadForm">
                <el-form-item :label="$t('auto.BatchImport.a015434e')" required>
                  <el-select v-model="uploadForm.product_id" :placeholder="$t('auto.BatchImport.cbf1cfff')" style="width: 100%" @change="loadPreview">
                    <el-option v-for="p in productList" :key="p.id" :label="p.name" :value="p.id" />
                  </el-select>
                </el-form-item>
                <el-form-item :label="$t('auto.BatchImport.e3127cc1')">
                  <el-radio-group v-model="uploadForm.format">
                    <el-radio-button label="auto">{{ $t('auto.BatchImport.0375083d') }}</el-radio-button>
                    <el-radio-button label="csv">CSV</el-radio-button>
                    <el-radio-button label="json">JSON</el-radio-button>
                  </el-radio-group>
                </el-form-item>
                <el-form-item :label="$t('auto.BatchImport.2a0c4740')" required>
                  <el-upload
                    ref="uploadRef"
                    :auto-upload="false"
                    :limit="1"
                    :on-change="handleFileChange"
                    :on-exceed="handleExceed"
                    accept=".csv,.json,.txt"
                    drag
                  >
                    <el-icon class="el-icon--upload"><upload-filled /></el-icon>
                    <div class="el-upload__text">{{ $t('auto.BatchImport.d3a5c7ba') }}<em>{{ $t('auto.BatchImport.2c808b02') }}</em></div>
                    <template #tip>
                      <div class="el-upload__tip">{{ $t('auto.BatchImport.85fb5725') }}</div>
                    </template>
                  </el-upload>
                </el-form-item>
              </el-form>
              <el-button type="primary" :icon="View" :disabled="!uploadForm.file" @click="loadPreview" style="margin-top: 8px">{{ $t('auto.BatchImport.3e2cae44') }}</el-button>
            </el-tab-pane>

            
            <el-tab-pane :label="$t('auto.BatchImport.3b658886')" name="paste">
              <el-form label-width="100px">
                <el-form-item :label="$t('auto.BatchImport.a015434e')" required>
                  <el-select v-model="pasteForm.product_id" :placeholder="$t('auto.BatchImport.cbf1cfff')" style="width: 100%">
                    <el-option v-for="p in productList" :key="p.id" :label="p.name" :value="p.id" />
                  </el-select>
                </el-form-item>
                <el-form-item :label="$t('auto.BatchImport.f51ef6d7')" required>
                  <el-input
                    v-model="pasteForm.jsonText"
                    type="textarea"
                    :rows="12"
                    placeholder='示例:&#10;[&#10;  {"title": "如何退货", "content": "请在订单页面申请", "category": "售后", "tags": ["退货","流程"]},&#10;  {"title": "产品 A 介绍", "content": "..."}&#10;]'
                  />
                </el-form-item>
                <el-form-item>
                  <el-button @click="loadPastePreview" :icon="View">{{ $t('auto.BatchImport.80510e5e') }}</el-button>
                  <el-button @click="loadPasteTemplate">{{ $t('auto.BatchImport.5806bb43') }}</el-button>
                </el-form-item>
              </el-form>
            </el-tab-pane>
          </el-tabs>
        </el-card>
      </el-col>

      
      <el-col :span="10">
        <el-card>
          <template #header>
            <div class="card-header">
              <span>{{ $t('auto.BatchImport.80510e5e') }}</span>
              <el-tag v-if="preview.length" type="info">共 {{ preview.length }} 条</el-tag>
            </div>
          </template>
          <el-table :data="preview" max-height="320" size="small" :empty-text="$t('auto.BatchImport.b53d71e4')">
            <el-table-column prop="title" :label="$t('auto.BatchImport.32c65d8d')" min-width="120" show-overflow-tooltip />
            <el-table-column prop="category" :label="$t('auto.BatchImport.d0771a42')" width="80" />
            <el-table-column :label="$t('auto.BatchImport.d9336549')" min-width="160" show-overflow-tooltip>
              <template #default="{ row }">
                <span class="content-prev">{{ truncateText(row.content, 60) }}</span>
              </template>
            </el-table-column>
            <el-table-column :label="$t('auto.BatchImport.14d34236')" width="120" show-overflow-tooltip>
              <template #default="{ row }">
                <span v-if="row.tags && row.tags.length">{{ row.tags.join(', ') }}</span>
                <span v-else>-</span>
              </template>
            </el-table-column>
          </el-table>
          <div style="margin-top: 12px; text-align: right">
            <el-button type="primary" :icon="Upload" :loading="importing" :disabled="!canImport" @click="handleBatchImport">
              确认导入 ({{ preview.length }})
            </el-button>
          </div>
        </el-card>
      </el-col>
    </el-row>

    
    <el-card v-if="lastResult" style="margin-top: 16px">
      <template #header>
        <div class="card-header">
          <span>
            <el-icon style="vertical-align: middle"><CircleCheck /></el-icon>
            导入完成
            <el-text type="info" size="small" style="margin-left: 8px">批次号: {{ lastResult.batch_no }}</el-text>
          </span>
          <el-button link type="primary" @click="lastResult = null">{{ $t('auto.BatchImport.b15d9127') }}</el-button>
        </div>
      </template>
      <el-row :gutter="16">
        <el-col :span="6">
          <el-statistic :title="$t('auto.BatchImport.9ed7d3ad')" :value="lastResult.total" />
        </el-col>
        <el-col :span="6">
          <el-statistic :title="$t('auto.BatchImport.330363df')" :value="lastResult.accepted" :value-style="{ color: '#10B981' }" />
        </el-col>
        <el-col :span="6">
          <el-statistic :title="$t('auto.BatchImport.acd5cb84')" :value="lastResult.rejected" :value-style="{ color: '#EF4444' }" />
        </el-col>
        <el-col :span="6">
          <el-statistic :title="$t('auto.BatchImport.b664352f')" :value="successRate" suffix="%" :precision="1" />
        </el-col>
      </el-row>
      <el-divider v-if="lastResult.errors && lastResult.errors.length" content-position="left">{{ $t('auto.BatchImport.6060a344') }}</el-divider>
      <el-table v-if="lastResult.errors && lastResult.errors.length" :data="lastResult.errors" max-height="200" size="small" border>
        <el-table-column type="index" :label="$t('auto.BatchImport.faaadc44')" width="60" />
        <el-table-column prop="value" :label="$t('auto.BatchImport.4604d502')" />
      </el-table>
    </el-card>

    
    <el-card style="margin-top: 16px">
      <template #header>
        <div class="card-header">
          <span>{{ $t('auto.BatchImport.5b86245c') }}</span>
        </div>
      </template>
      <pre class="csv-help">
title,content,category,tags
"如何退货","请在订单页面申请",售后,"退货,流程"
"产品 A 介绍","产品 A 是我们的明星产品",产品,"产品A,介绍"</pre>
      <p style="color: #909399; font-size: 12px">
        第一行为表头,必须包含 <code>content</code> 列(标题、分类、标签可选)。
        多标签用英文逗号分隔。
      </p>
    </el-card>
  </div>
</template>

<script setup>
import i18n from '@/i18n'

import { ref, reactive, computed, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Upload, View, CircleCheck, UploadFilled } from '@element-plus/icons-vue'
import { knowledgeMerchantAPI } from '@/api/knowledgeMerchant'
import { ragProductConfigAPI } from '@/api/ragProductConfig'

const { t } = useI18n()

const activeTab = ref('upload')
const uploading = ref(false)
const importing = ref(false)

const productList = ref([])
const preview = ref([])
const lastResult = ref(null)

const uploadForm = reactive({
  product_id: '',
  format: 'auto',
  file: null
})

const pasteForm = reactive({
  product_id: '',
  jsonText: ''
})

const canImport = computed(() => preview.value.length > 0 && !importing.value)
const successRate = computed(() => {
  if (!lastResult.value || !lastResult.value.total) return 0
  return ((lastResult.value.accepted / lastResult.value.total) * 100).toFixed(1)
})

const loadProducts = async () => {
  try {
    const res = await ragProductConfigAPI.listProducts()
    if (Array.isArray(res)) productList.value = res
    else if (res?.items) productList.value = res.items
  } catch (e) {
    console.error(t('auto.BatchImport.4c46f3b1'), e)
  }
}

const handleTabChange = (tab) => {
  preview.value = []
  lastResult.value = null
}

const handleFileChange = (file) => {
  if (file.size > 20 * 1024 * 1024) {
    ElMessage.error(i18n.global.t(t('auto.BatchImport.55b3530f')))
    return false
  }
  uploadForm.file = file.raw
  return true
}

const handleExceed = () => {
  ElMessage.warning(i18n.global.t(t('auto.BatchImport.c4e4c6d3')))
}

const loadPreview = () => {
  if (!uploadForm.file) {
    ElMessage.warning(i18n.global.t(t('auto.BatchImport.3c2412f9')))
    return
  }
  const reader = new FileReader()
  reader.onload = (e) => {
    const text = e.target.result
    parseTextToPreview(text, uploadForm.format === 'auto' ? detectFormat(uploadForm.file.name) : uploadForm.format)
  }
  reader.readAsText(uploadForm.file)
};

const loadPastePreview = () => {
  if (!pasteForm.jsonText.trim()) {
    ElMessage.warning(i18n.global.t(t('auto.BatchImport.1ae7de68')))
    return
  }
  parseTextToPreview(pasteForm.jsonText, 'json')
}

const detectFormat = (filename) => {
  const lower = filename.toLowerCase()
  if (lower.endsWith('.json')) return 'json'
  if (lower.endsWith('.csv')) return 'csv'
  return 'json'
}

const parseTextToPreview = (text, format) => {
  try {
    if (format === 'csv') {
      preview.value = parseCSVClient(text)
    } else {
      preview.value = parseJSONClient(text)
    }
    if (preview.value.length === 0) {
      ElMessage.warning(i18n.global.t(t('auto.BatchImport.3dacfab7')))
    } else {
      ElMessage.success(`解析成功,共 ${preview.value.length} 条`)
    }
  } catch (e) {
    preview.value = []
    ElMessage.error(t('auto.BatchImport.fa5e42f4') + e.message)
  }
}

const parseCSVClient = (text) => {
  const lines = text.split(/\r?\n/).filter(l => l.trim())
  if (lines.length < 2) throw new Error('CSV 必须至少包含表头和一行数据')
  const headers = parseCSVLine(lines[0])
  const idxContent = headers.findIndex(h => /content|内容|text|body/i.test(h))
  if (idxContent < 0) throw new Error('CSV 缺少 content 列')
  const idxTitle = headers.findIndex(h => /title|标题|name/i.test(h))
  const idxCategory = headers.findIndex(h => /category|分类/i.test(h))
  const idxTags = headers.findIndex(h => /tags|标签/i.test(h))
  const out = []
  for (let i = 1; i < lines.length; i++) {
    const cols = parseCSVLine(lines[i])
    if (!cols[idxContent]) continue
    out.push({
      title: idxTitle >= 0 ? (cols[idxTitle] || '') : '',
      content: cols[idxContent],
      category: idxCategory >= 0 ? (cols[idxCategory] || '') : '',
      tags: idxTags >= 0 ? (cols[idxTags] || '').split(',').map(s => s.trim()).filter(Boolean) : []
    })
  }
  return out
}

const parseCSVLine = (line) => {
  const out = []
  let cur = ''
  let inQuote = false
  for (let i = 0; i < line.length; i++) {
    const c = line[i]
    if (c === '"') {
      if (inQuote && line[i + 1] === '"') { cur += '"'; i++ }
      else inQuote = !inQuote
    } else if (c === ',' && !inQuote) {
      out.push(cur)
      cur = ''
    } else {
      cur += c
    }
  }
  out.push(cur)
  return out
}

const parseJSONClient = (text) => {
  const data = JSON.parse(text)
  if (Array.isArray(data)) return data
  if (data.items) return data.items
  if (data.data) return data.data
  if (data.list) return data.list
  if (data.documents) return data.documents
  throw new Error('JSON 结构不匹配,需为数组或包含 items/data/list/documents 字段')
}

const loadPasteTemplate = () => {
  pasteForm.jsonText = JSON.stringify([
    { title: i18n.global.t(t('auto.BatchImport.9570d089')), content: i18n.global.t(t('auto.BatchImport.d6ba401b')), category: 'FAQ', tags: [t('auto.BatchImport.1a63ac23'), 'FAQ'] },
    { title: i18n.global.t(t('auto.BatchImport.b8d7eb92')), content: i18n.global.t(t('auto.BatchImport.4e27697d')), category: t('auto.BatchImport.a015434e'), tags: [t('auto.BatchImport.61a3ec66')] }
  ], null, 2)
}

const handleBatchImport = async () => {
  if (!preview.value.length) return
  try {
    await ElMessageBox.confirm(
      `确认导入 ${preview.value.length} 条记录吗?`,
      t('auto.BatchImport.407f3df2'),
      { type: 'info' }
    )
  } catch { return }

  importing.value = true
  try {
    let res
    if (activeTab.value === 'upload') {
      const fd = new FormData()
      fd.append('product_id', uploadForm.product_id)
      fd.append('format', uploadForm.format)
      fd.append('file', uploadForm.file)
      res = await knowledgeMerchantAPI.batchUpload(fd)
    } else {
      res = await knowledgeMerchantAPI.batchImport({
        product_id: pasteForm.product_id,
        items: preview.value
      })
    }
    if (res) {
      lastResult.value = res
      ElMessage.success(`导入完成:成功 ${res.accepted} / 失败 ${res.rejected}`)
      preview.value = [];
      if (activeTab.value === 'upload') uploadForm.file = null
      else pasteForm.jsonText = ''
    }
  } catch (e) {
    ElMessage.error(t('auto.BatchImport.45332d13') + (e.message || ''))
  } finally {
    importing.value = false
  }
}

const truncateText = (s, n) => s ? (s.length > n ? s.substring(0, n) + '...' : s) : ''

onMounted(() => {
  loadProducts()
})
</script>

<style scoped lang="scss">
.batch-import-page {
  padding: 0;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.content-prev {
  font-size: 12px;
  color: #606266;
}

.csv-help {
  background: #fafafa;
  border: 1px solid #ebeef5;
  border-radius: 4px;
  padding: 12px;
  font-size: 12px;
  font-family: monospace;
  color: #303133;
  white-space: pre-wrap;
  word-break: break-all;
}
</style>
