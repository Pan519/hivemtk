<template>
  <div class="tag-segmentation-page">
    <el-card class="header-card">
      <div class="header-content">
        <h2>{{ $t('标签分层') }}</h2>
        <p class="subtitle">{{ $t('auto.List.41.996ed7') }}</p>
      </div>
      <div class="header-actions">
        <el-button type="primary" @click="showTagDialog()">
          <el-icon><Plus /></el-icon>
          {{ $t('新增标签') }}
        </el-button>
        <el-button @click="refreshAll">
          <el-icon><Refresh /></el-icon>
          {{ $t('刷新') }}
        </el-button>
      </div>
    </el-card>

    <el-tabs v-model="activeTab" class="content-tabs">
      <el-tab-pane ::label="$t('auto.List.1.7114ff')" name="tags">
        <div class="filter-bar">
          <el-input v-model="tagSearch" ::placeholder="$t('auto.List.2.24eabe')" clearable style="width: 220px" />
          <el-select v-model="tagTypeFilter" ::placeholder="$t('auto.List.3.f19b72')" clearable style="width: 160px">
            <el-option ::label="$t('auto.List.4.1b14ee')" value="manual" />
            <el-option ::label="$t('auto.List.5.c08f63')" value="auto" />
            <el-option ::label="$t('auto.List.6.b67742')" value="system" />
          </el-select>
        </div>
        <el-table :data="filteredTags" v-loading="loading.tags" stripe>
          <template #empty><el-empty :description="$t('auto.List.7.dafa8c')" /></template>
          <el-table-column prop="name" :label="$t('auto.List.8.341fe8')" min-width="180" show-overflow-tooltip />
          <el-table-column :label="$t('auto.List.9.226b09')" min-width="90">
            <template #default="{ row }">
              <el-tag :type="getTypeColor(row.tag_type || row.type || 'manual')" size="small">
                {{ getTypeText(row.tag_type || row.type || 'manual') }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column :label="$t('auto.List.10.d0771a')" min-width="120" show-overflow-tooltip>
            <template #default="{ row }">{{ row.group || row.category || '默认' }}</template>
          </el-table-column>
          <el-table-column :label="$t('auto.List.11.71cb56')" min-width="90" align="center">
            <template #default="{ row }">{{ row.user_count ?? row.userCount ?? 0 }}</template>
          </el-table-column>
          <el-table-column :label="$t('auto.List.12.3fea7c')" min-width="90" align="center">
            <template #default="{ row }">
              <el-tag :type="row.is_active ? 'success' : 'info'" size="small">
                {{ row.is_active ? '启用' : '禁用' }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column :label="$t('auto.List.13.eca37c')" min-width="160" show-overflow-tooltip>
            <template #default="{ row }">{{ formatTime(row.created_at || row.createdAt) }}</template>
          </el-table-column>
          <el-table-column :label="$t('auto.List.14.2b6bc0')" min-width="160" fixed="right">
            <template #default="{ row }">
              <el-button link type="primary" size="small" @click="showTagDialog(row)">{{ $t('auto.List.42.95b351') }}</el-button>
              <el-button link type="danger" size="small" @click="deleteTag(row)">{{ $t('auto.List.43.2f4aad') }}</el-button>
            </template>
          </el-table-column>
        </el-table>
      </el-tab-pane>

      <el-tab-pane :label="$t('auto.List.15.ad4676')" name="rules">
        <el-card>
          <template #header>
            <div class="card-header">
              <span>{{ $t('auto.List.44.c7a994') }}</span>
              <el-button type="primary" size="small" @click="showRuleDialog()">
                <el-icon><Plus /></el-icon> 新增规则
              </el-button>
            </div>
          </template>
          <el-table :data="tagRuleList" v-loading="loading.rules" stripe>
            <template #empty><el-empty :description="$t('auto.List.16.60639b')" /></template>
            <el-table-column prop="name" :label="$t('auto.List.17.870802')" min-width="150" />
            <el-table-column prop="condition" :label="$t('auto.List.18.65939c')" min-width="260" show-overflow-tooltip />
            <el-table-column prop="tagName" :label="$t('auto.List.19.fedbca')" width="140">
              <template #default="{ row }">
                <el-tag size="small">{{ row.tagName }}</el-tag>
              </template>
            </el-table-column>
            <el-table-column prop="status" :label="$t('auto.List.12.3fea7c')" width="100">
              <template #default="{ row }">
                <el-tag :type="getStatusTagType(row.status)" size="small">
                  {{ getStatusLabel(row.status) }}
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column prop="matchedCount" :label="$t('auto.List.20.2f0aac')" width="100" align="center" />
            <el-table-column :label="$t('auto.List.14.2b6bc0')" width="180" fixed="right">
              <template #default="{ row }">
                <el-button link type="primary" @click="showRuleDialog(row)">{{ $t('auto.List.42.95b351') }}</el-button>
                <el-button link type="danger" @click="deleteRule(row)">{{ $t('auto.List.43.2f4aad') }}</el-button>
              </template>
            </el-table-column>
          </el-table>
        </el-card>
      </el-tab-pane>

      <el-tab-pane :label="$t('auto.List.21.b61fbd')" name="strategy">
        <el-card v-loading="loading.strategy">
          <template #header>
            <div class="card-header">
              <span>{{ $t('auto.List.45.e9f023') }}</span>
              <el-button type="primary" size="small" @click="saveStrategy">{{ $t('auto.List.46.2d3f52') }}</el-button>
            </div>
          </template>
          <el-table :data="strategyList" stripe>
            <template #empty><el-empty :description="$t('auto.List.22.f4c3a4')" /></template>
            <el-table-column prop="layer" :label="$t('auto.List.23.a8d7cf')" width="160">
              <template #default="{ row }">
                <el-tag :type="getLayerColor(row.layer)">{{ row.layer }}</el-tag>
              </template>
            </el-table-column>
            <el-table-column prop="description" :label="$t('auto.List.24.dcced8')" min-width="200" />
            <el-table-column :label="$t('auto.List.25.c0eb17')" min-width="280">
              <template #default="{ row }">
                <el-select v-model="row.tagIds" multiple style="width: 100%" :placeholder="$t('auto.List.26.f9b05d')">
                  <el-option v-for="t in tags" :key="t.id" :label="t.name" :value="t.id" />
                </el-select>
              </template>
            </el-table-column>
          </el-table>
        </el-card>
      </el-tab-pane>

      <el-tab-pane :label="$t('auto.List.27.acc4bd')" name="stats">
        <el-row :gutter="20" class="stat-row">
          <el-col :span="6">
            <el-card class="stat-card">
              <div class="stat-label">{{ $t('auto.List.47.a7f020') }}</div>
              <div class="stat-value">{{ stats.total || 0 }}</div>
            </el-card>
          </el-col>
          <el-col :span="6">
            <el-card class="stat-card">
              <div class="stat-label">{{ $t('auto.List.48.87ae72') }}</div>
              <div class="stat-value" style="color: #4F46E5">{{ stats.auto || 0 }}</div>
            </el-card>
          </el-col>
          <el-col :span="6">
            <el-card class="stat-card">
              <div class="stat-label">{{ $t('auto.List.49.68004d') }}</div>
              <div class="stat-value" style="color: #10B981">{{ stats.manual || 0 }}</div>
            </el-card>
          </el-col>
          <el-col :span="6">
            <el-card class="stat-card">
              <div class="stat-label">{{ $t('auto.List.50.11d835') }}</div>
              <div class="stat-value" style="color: #F59E0B">{{ stats.matchedUsers || 0 }}</div>
            </el-card>
          </el-col>
        </el-row>
        <el-card>
          <template #header><span>{{ $t('auto.List.51.55ebe6') }}</span></template>
          <el-table :data="stats.topUsed" v-loading="loading.stats" stripe>
            <template #empty><el-empty :description="$t('auto.List.28.a3b437')" /></template>
            <el-table-column type="index" :label="$t('auto.List.29.a4dc00')" width="80" />
            <el-table-column prop="name" :label="$t('auto.List.30.14d342')" min-width="160" />
            <el-table-column prop="userCount" :label="$t('auto.List.11.71cb56')" width="120" align="center" />
            <el-table-column prop="ratio" :label="$t('auto.List.31.c909e6')" width="200">
              <template #default="{ row }">
                <el-progress :percentage="Number(row.ratio || 0)" :stroke-width="8" />
              </template>
            </el-table-column>
          </el-table>
        </el-card>
      </el-tab-pane>
    </el-tabs>

    <el-dialog v-model="tagDialogVisible" :title="tagDialogTitle" width="560px">
      <el-form :model="tagForm" :rules="tagFormRules" ref="tagFormRef" label-width="100px">
        <el-form-item :label="$t('auto.List.8.341fe8')" prop="name">
          <el-input v-model="tagForm.name" :placeholder="$t('auto.List.32.6f81f3')" />
        </el-form-item>
        <el-form-item :label="$t('auto.List.33.200948')" prop="type">
          <el-select v-model="tagForm.type" style="width: 100%">
            <el-option :label="$t('auto.List.34.2a3e7f')" value="manual" />
            <el-option :label="$t('auto.List.35.3aed2c')" value="auto" />
            <el-option :label="$t('auto.List.36.8a8b89')" value="system" />
          </el-select>
        </el-form-item>
        <el-form-item :label="$t('auto.List.10.d0771a')" prop="category">
          <el-input v-model="tagForm.category" :placeholder="$t('auto.List.37.d34f6c')" />
        </el-form-item>
        <el-form-item :label="$t('auto.List.38.3bdd08')">
          <el-input v-model="tagForm.description" type="textarea" :rows="3" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="tagDialogVisible = false">{{ $t('auto.List.52.625fb2') }}</el-button>
        <el-button type="primary" @click="submitTag">{{ $t('auto.List.53.38cf16') }}</el-button>
      </template>
    </el-dialog>

    <el-dialog v-model="ruleDialogVisible" :title="ruleDialogTitle" width="640px">
      <el-form :model="ruleForm" :rules="ruleFormRules" ref="ruleFormRef" label-width="100px">
        <el-form-item :label="$t('auto.List.17.870802')" prop="name">
          <el-input v-model="ruleForm.name" :placeholder="$t('auto.List.39.40f3f1')" />
        </el-form-item>
        <el-form-item :label="$t('auto.List.18.65939c')" prop="condition">
          <el-input v-model="ruleForm.condition" type="textarea" :rows="3" :placeholder="$t('auto.List.40.7ad15c')" />
        </el-form-item>
        <el-form-item :label="$t('auto.List.19.fedbca')" prop="tagId">
          <el-select v-model="ruleForm.tagId" filterable :placeholder="$t('auto.List.26.f9b05d')" style="width: 100%">
            <el-option v-for="t in tags" :key="t.id" :label="t.name" :value="t.id" />
          </el-select>
        </el-form-item>
        <el-form-item :label="$t('auto.List.12.3fea7c')">
          <el-switch v-model="ruleForm.status" active-value="enabled" inactive-value="disabled" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="ruleDialogVisible = false">{{ $t('auto.List.52.625fb2') }}</el-button>
        <el-button type="primary" @click="submitRule">{{ $t('auto.List.53.38cf16') }}</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import i18n from '@/i18n'

import { ref, computed, reactive, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Plus, Refresh } from '@element-plus/icons-vue'
import { TagSegmentationApi } from '@/api/tagSegmentation.js'
import { getEnabledLabel, getEnabledTagType } from '@/constants/enabled';

const getStatusLabel = (s) => getEnabledLabel(s)
const getStatusTagType = (s) => getEnabledTagType(s)

const { t } = useI18n()

const activeTab = ref('tags')
const loading = reactive({ tags: false, rules: false, strategy: false, stats: false })

const tags = ref([])
const tagRuleList = ref([])
const strategyList = ref([])
const stats = ref({ topUsed: [] })

const tagSearch = ref('')
const tagTypeFilter = ref('')

const tagDialogVisible = ref(false)
const tagDialogTitle = ref(t('auto.List.54.a1d3bf'))
const tagFormRef = ref()
const tagForm = ref({ id: 0, name: '', type: 'manual', category: '', description: '' })
const tagFormRules = {
  name: [{ required: true, message: i18n.global.t(t('auto.List.32.6f81f3')), trigger: 'blur' }],
  type: [{ required: true, message: i18n.global.t(t('auto.List.55.f5d604')), trigger: 'change' }]
}

const ruleDialogVisible = ref(false)
const ruleDialogTitle = ref(t('auto.List.56.ba3c80'))
const ruleFormRef = ref()
const ruleForm = ref({ id: 0, name: '', condition: '', tagId: '', status: 'enabled' })
const ruleFormRules = {
  name: [{ required: true, message: i18n.global.t(t('auto.List.39.40f3f1')), trigger: 'blur' }],
  condition: [{ required: true, message: i18n.global.t(t('auto.List.57.40c86c')), trigger: 'blur' }],
  tagId: [{ required: true, message: i18n.global.t(t('auto.List.58.ce7cf5')), trigger: 'change' }]
}

const filteredTags = computed(() => {
  return tags.value.filter(t => {
    const matchName = !tagSearch.value || t.name?.includes(tagSearch.value)
    const matchType = !tagTypeFilter.value || t.type === tagTypeFilter.value
    return matchName && matchType
  })
})

const getTypeColor = (type) => {
  const map = { manual: 'primary', auto: 'success', system: 'info' }
  return map[type]}
const getTypeText = (type) => {
  const map = { manual: t('auto.List.34.2a3e7f'), auto: t('auto.List.35.3aed2c'), system: t('auto.List.36.8a8b89') }
  return map[type] || type
}
const getLayerColor = (layer) => {
  if (!layer || typeof layer !== 'string') return 'primary'
  if (layer.includes(t('auto.List.59.2bafca')) || layer.includes(t('auto.List.60.fc7e38'))) return 'success'
  if (layer.includes(t('auto.List.61.3adc53'))) return 'danger'
  if (layer.includes(t('auto.List.62.5f0c0b'))) return 'warning'
  return 'primary'
}

const loadTags = async () => {
  loading.tags = true
  try {
    const res = await TagSegmentationApi.getTags()
    const data = res
    tags.value = Array.isArray(data) ? data : (data?.items || data?.list || [])
  } catch (e) {
    tags.value = []
  } finally {
    loading.tags = false
  }
}

const loadRules = async () => {
  loading.rules = true
  try {
    const res = await TagSegmentationApi.getTagRules()
    const data = res
    tagRuleList.value = Array.isArray(data) ? data : (data?.items || data?.list || [])
  } catch (e) {
    tagRuleList.value = []
  } finally {
    loading.rules = false
  }
}

const loadStrategy = async () => {
  loading.strategy = true
  try {
    const res = await TagSegmentationApi.getLayerStrategy()
    const data = res
    strategyList.value = Array.isArray(data) ? data : (data?.items || data?.list || [])
  } catch (e) {
    strategyList.value = []
  } finally {
    loading.strategy = false
  }
}

const loadStats = async () => {
  loading.stats = true
  try {
    const res= await TagSegmentationApi.getTagStats()
    const d = res || {}
    stats.value = { topUsed: Array.isArray(d.topUsed) ? d.topUsed : (Array.isArray(d) ? d : []) }
  } catch (e) {
    stats.value = {}
  } finally {
    loading.stats = false
  }
}

const refreshAll = () => {
  loadTags()
  loadRules()
  loadStrategy()
  loadStats()
}

const showTagDialog = (row) => {
  if (row) {
    tagForm.value = { ...row }
    tagDialogTitle.value = t('auto.List.63.4f3b1f')
  } else {
    tagForm.value = { id: 0, name: '', code: '', type: 'manual', category: '', description: '' }
    tagDialogTitle.value = t('auto.List.54.a1d3bf')
  }
  tagDialogVisible.value = true
}

const formatDate = (date) => {
  if (!date) return ''
  const d = new Date(date)
  const year = d.getFullYear()
  const month = (d.getMonth() + 1).toString().padStart(2, '0')
  const day = d.getDate().toString().padStart(2, '0')
  return `${year}-${month}-${day}`
}

const formatTime = (timeStr) => {
  if (!timeStr) return ''
  const date = new Date(timeStr)
  const now = new Date()
  const diff = now.getTime() - date.getTime()
  const days = Math.floor(diff / (1000 * 60 * 60 * 24))

  if (days === 0) {
    const hours = Math.floor(diff / (1000 * 60 * 60))
    if (hours === 0) {
      const minutes = Math.floor(diff / (1000 * 60))
      return minutes <= 1 ? t('auto.List.64.4181f7') : `${minutes}分钟前`
    }
    return `${hours}小时前`
  } else if (days === 1) {
    return t('auto.List.65.2f8d6f')
  } else if (days < 7) {
    return `${days}天前`
  } else {
    return formatDate(date)
  }
}

const genTagCode = (name) => {
  const base = (name || '').toString().trim()
    .toLowerCase()
    .replace(/[^a-z0-9\u4e00-\u9fa5]+/g, '_')
    .replace(/^_+|_+$/g, '')
  return base ? `tag_${base}` : `tag_${Date.now()}`
}

const submitTag = async () => {
  if (!tagFormRef.value) return
  await tagFormRef.value.validate(async (valid) => {
    if (!valid) return
    try {
      const payload = {
        name: tagForm.value.name,
        code: tagForm.value.code || genTagCode(tagForm.value.name),
        group: tagForm.value.category || 'default',
        description: tagForm.value.description || ''
      };
      if (tagForm.value.id) {
        await TagSegmentationApi.updateTags({ id: tagForm.value.id, ...payload })
      } else {
        await TagSegmentationApi.createTag(payload)
      }
      ElMessage.success(tagForm.value.id ? t('auto.List.66.55aa63') : t('auto.List.67.a5bfd7'))
      tagDialogVisible.value = false
      loadTags()
    } catch (e) {
      ElMessage.error(i18n.global.t(t('auto.List.68.6de920')))
    }
  })
}

const deleteTag = async (row) => {
  try {
    await ElMessageBox.confirm(`确定删除标签 "${row.name}" 吗？`, t('auto.List.69.e83a25'), { type: 'warning' })
    await TagSegmentationApi.deleteTag(row.id)
    ElMessage.success(i18n.global.t(t('auto.List.70.0007d1')))
    loadTags()
  } catch (e) {
    if (e !== 'cancel') ElMessage.error(i18n.global.t(t('auto.List.71.acf066')))
  }
}

const showRuleDialog = (row) => {
  if (row) {
    ruleForm.value = { ...row }
    ruleDialogTitle.value = t('auto.List.72.a1d940')
  } else {
    ruleForm.value = { id: 0, name: '', condition: '', tagId: '', status: 'enabled' }
    ruleDialogTitle.value = t('auto.List.56.ba3c80')
  }
  ruleDialogVisible.value = true
}

const submitRule = async () => {
  if (!ruleFormRef.value) return
  await ruleFormRef.value.validate(async (valid) => {
    if (!valid) return
    const tagName = tags.value.find(t => t.id === ruleForm.value.tagId)?.name
    const payload = { ...ruleForm.value, tagName }
    try {
      if (ruleForm.value.id) {
        await TagSegmentationApi.updateTagRule(ruleForm.value.id, payload)
      } else {
        await TagSegmentationApi.saveTagRule(payload)
      }
      ElMessage.success(ruleForm.value.id ? t('auto.List.66.55aa63') : t('auto.List.67.a5bfd7'))
      ruleDialogVisible.value = false
      loadRules()
    } catch (e) {
      ElMessage.error(i18n.global.t(t('auto.List.68.6de920')))
    }
  })
}

const deleteRule = async (row) => {
  try {
    await ElMessageBox.confirm(`确定删除规则 "${row.name}" 吗？`, t('auto.List.69.e83a25'), { type: 'warning' })
    await TagSegmentationApi.deleteTagRule(row.id)
    ElMessage.success(i18n.global.t(t('auto.List.70.0007d1')))
    loadRules()
  } catch (e) {
    if (e !== 'cancel') ElMessage.error(i18n.global.t(t('auto.List.71.acf066')))
  }
}

const saveStrategy = async () => {
  try {
    const payload = strategyList.value.map(r => ({
      layer: r.layer,
      description: r.description || ''
    }));
    await TagSegmentationApi.saveLayerStrategy(payload)
    ElMessage.success(i18n.global.t(t('auto.List.73.2864c2')))
  } catch (e) {
    ElMessage.error(i18n.global.t(t('auto.List.68.6de920')))
  }
}

onMounted(() => {
  refreshAll()
})
</script>

<style scoped lang="scss">
.tag-segmentation-page { padding: 20px; }
.header-card {
  margin-bottom: 20px;
  :deep(.el-card__body) {
    display: flex;
    justify-content: space-between;
    align-items: center;
  }
  .header-content h2 { margin: 0 0 8px 0; }
  .subtitle { color: #909399; margin: 0; }
  .header-actions { display: flex; gap: 10px; }
}
.content-tabs { background: #fff; padding: 16px; border-radius: 4px; }
.filter-bar { display: flex; gap: 12px; margin-bottom: 16px; }
.card-header { display: flex; justify-content: space-between; align-items: center; }
.stat-row { margin-bottom: 20px; }
.stat-card {
  text-align: center;
  .stat-label { color: #909399; font-size: 14px; margin-bottom: 10px; }
  .stat-value { font-size: 28px; font-weight: bold; }
}
</style>
