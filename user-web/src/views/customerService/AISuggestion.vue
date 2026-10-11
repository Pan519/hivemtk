<template>
  <div class="ai-suggestion-page">
    <el-card class="header-card">
      <div>
        <h2>{{ $t('auto.AISuggestion.2.5d4180') }}</h2>
        <p class="subtitle">{{ $t('auto.AISuggestion.9.bc224d') }}</p>
      </div>
    </el-card>

    <el-card>
      <template #header>
        <div class="card-header">
          <el-space>
            <el-input
              v-model="sessionId"
              ::placeholder="$t('auto.AISuggestion.1.42dd8c')"
              style="width: 280px"
              clearable
              @keyup.enter="loadSuggestions"
            />
            <el-button type="primary" @click="loadSuggestions">
              <el-icon><Search /></el-icon>
              {{ $t('加载') }}
            </el-button>
          </el-space>
        </div>
      </template>

      <el-table :data="suggestions" v-loading="loading" stripe>
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="suggestion" :label="$t('auto.AISuggestion.2.5d4180')" min-width="320" show-overflow-tooltip />
        <el-table-column prop="confidence" :label="$t('auto.AISuggestion.3.8f3599')" width="120" align="center">
          <template #default="{ row }">
            <el-progress
              :percentage="Math.round((row.confidence || 0) * 100)"
              :stroke-width="10"
              :color="confidenceColor(row.confidence)"
            />
          </template>
        </el-table-column>
        <el-table-column prop="source" :label="$t('auto.AISuggestion.4.26ca20')" width="120">
          <template #default="{ row }">
            <el-tag size="small">{{ row.source || 'rule' }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="is_used" :label="$t('auto.AISuggestion.5.453597')" width="120" align="center">
          <template #default="{ row }">
            <el-tag v-if="row.is_used" type="success">{{ $t('auto.AISuggestion.10.5dadca') }}</el-tag>
            <el-tag v-else type="info">{{ $t('auto.AISuggestion.11.d05325') }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="created_at" :label="$t('auto.AISuggestion.6.b683c8')" min-width="160">
          <template #default="{ row }">
            {{ formatTime(row.created_at) }}
          </template>
        </el-table-column>
        <el-table-column :label="$t('auto.AISuggestion.7.2b6bc0')" width="120" fixed="right">
          <template #default="{ row }">
            <el-button
              v-if="!row.is_used"
              link
              type="success"
              @click="handleUse(row)"
            >{{ $t('auto.AISuggestion.12.6cffff') }}</el-button>
            <span v-else style="color: #909399">{{ $t('auto.AISuggestion.13.b59b00') }}</span>
          </template>
        </el-table-column>
        <template #empty>
          <el-empty ::description="$t('auto.AISuggestion.8.6f28b6')" />
        </template>
      </el-table>
    </el-card>
  </div>
</template>

<script setup>
import i18n from '@/i18n'

import { ref, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { ElMessage } from 'element-plus'
import { Search } from '@element-plus/icons-vue'
import { getAISuggestions, useAISuggestion } from '@/api/customerService.js'

const { t } = useI18n()

const loading = ref(false)
const sessionId = ref('')
const suggestions = ref([])

const formatTime = (val) => {
  if (!val) return '-'
  const d = new Date(val)
  if (isNaN(d.getTime())) return '-'
  return d.toLocaleString('zh-CN', { hour12: false })
}

const confidenceColor = (c) => {
  if (c >= 0.8) return '#10B981'
  if (c >= 0.5) return '#F59E0B'
  return '#EF4444'
}

const loadSuggestions = async () => {
  if (!sessionId.value) {
    ElMessage.warning(i18n.global.t(t('auto.AISuggestion.14.8b9019')))
    return
  }
  loading.value = true
  try {
    const res = await getAISuggestions(sessionId.value)
    const data = res || []
    suggestions.value = Array.isArray(data) ? data : data.list || []
  } catch (e) {
    ElMessage.error(i18n.global.t(t('auto.AISuggestion.15.6cacd0')))
    suggestions.value = []
  } finally {
    loading.value = false
  }
}

const handleUse = async (row) => {
  try {
    await useAISuggestion(row.id)
    ElMessage.success(i18n.global.t(t('auto.AISuggestion.10.5dadca')))
    row.is_used = true
  } catch (e) {
    ElMessage.error(i18n.global.t(t('auto.AISuggestion.16.b257f6')))
  }
}

onMounted(() => {})
</script>

<style scoped lang="scss">
.ai-suggestion-page { padding: 20px; }
.header-card {
  margin-bottom: 20px;
  :deep(.el-card__body) { display: flex; justify-content: space-between; align-items: center; }
  h2 { margin: 0 0 8px 0; }
  .subtitle { color: #909399; margin: 0; }
}
.card-header { display: flex; justify-content: space-between; align-items: center; }
</style>
