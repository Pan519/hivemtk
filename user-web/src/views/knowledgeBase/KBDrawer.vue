<template>
  <el-drawer
    :model-value="modelValue"
    @update:model-value="$emit('update:modelValue', $event)"
    :title="drawerTitle"
    size="720px"
    direction="rtl"
    :destroy-on-close="false"
  >
    <div v-loading="loading" class="kb-drawer">
      <template v-if="kbData">
        
        <el-card shadow="never" class="section">
          <template #header>
            <div class="section-header">
              <el-icon><InfoFilled /></el-icon>
              <span>{{ $t('auto.KBDrawer.9e5ffa06') }}</span>
            </div>
          </template>
          <el-descriptions :column="2" border>
            <el-descriptions-item label="ID">{{ kbData.id }}</el-descriptions-item>
            <el-descriptions-item :label="$t('auto.KBDrawer.031c076a')">{{ kbData.kb_code || '-' }}</el-descriptions-item>
            <el-descriptions-item :label="$t('auto.KBDrawer.d7ec2d3f')" :span="2">
              {{ kbData.name }}
            </el-descriptions-item>
            <el-descriptions-item :label="$t('auto.KBDrawer.226b0912')" :span="1">
              <el-tag :type="getTypeTagType(kbData.kb_type)" size="small">
                {{ getTypeLabel(kbData.kb_type) }}
              </el-tag>
            </el-descriptions-item>
            <el-descriptions-item :label="$t('auto.KBDrawer.3fea7ca7')" :span="1">
              <el-tag :type="kbData.status === 1 ? 'success' : 'info'" size="small">
                {{ kbData.status === 1 ? '启用' : '禁用' }}
              </el-tag>
            </el-descriptions-item>
            <el-descriptions-item :label="$t('auto.KBDrawer.3bdd08ad')" :span="2">
              {{ kbData.description || '-' }}
            </el-descriptions-item>
            <el-descriptions-item :label="$t('auto.KBDrawer.eca37cb0')">
              {{ formatTime(kbData.created_at) }}
            </el-descriptions-item>
            <el-descriptions-item :label="$t('auto.KBDrawer.a001a226')">
              {{ formatTime(kbData.updated_at) }}
            </el-descriptions-item>
          </el-descriptions>
        </el-card>

        
        <el-card shadow="never" class="section">
          <template #header>
            <div class="section-header">
              <el-icon><DataLine /></el-icon>
              <span>{{ $t('auto.KBDrawer.f56ae67a') }}</span>
            </div>
          </template>
          <el-row :gutter="16">
            <el-col :span="8">
              <div class="mini-stat">
                <div class="mini-stat-label">{{ statsLabel.item }}</div>
                <div class="mini-stat-value">{{ stats.item_count ?? '—' }}</div>
                <div class="mini-stat-hint">{{ $t('auto.KBDrawer.8f5c4bfd') }}</div>
              </div>
            </el-col>
            <el-col :span="8">
              <div class="mini-stat">
                <div class="mini-stat-label">{{ $t('auto.KBDrawer.9b4c8f90') }}</div>
                <div class="mini-stat-value text-primary">{{ stats.agent_count }}</div>
              </div>
            </el-col>
            <el-col :span="8">
              <div class="mini-stat">
                <div class="mini-stat-label">{{ $t('auto.KBDrawer.ba457909') }}</div>
                <div class="mini-stat-value text-success">{{ stats.hit_count ?? '—' }}</div>
                <div class="mini-stat-hint">{{ $t('auto.KBDrawer.8f5c4bfd') }}</div>
              </div>
            </el-col>
          </el-row>
        </el-card>

        
        <el-card shadow="never" class="section">
          <template #header>
            <div class="section-header">
              <el-icon><Connection /></el-icon>
              <span>{{ $t('auto.KBDrawer.ef41be05') }}</span>
              <el-button
                type="primary"
                link
                size="small"
                :loading="loadingAgents"
                @click="loadAgentBindings"
                style="margin-left: auto"
              >
                <el-icon><Refresh /></el-icon>刷新
              </el-button>
            </div>
          </template>
          <el-empty
            v-if="!loadingAgents && agentBindings.length === 0"
            :description="$t('auto.KBDrawer.e792d8b4')"
            :image-size="80"
          />
          <el-table v-else :data="agentBindings" stripe size="small" border>
            <el-table-column prop="id" label="ID" width="70" align="center" />
            <el-table-column :label="$t('auto.KBDrawer.63ecdffc')" min-width="160" show-overflow-tooltip>
              <template #default="{ row }">
                <span class="agent-link" role="button" tabindex="0" @click="goEditAgent(row)" @keydown.enter.prevent="goEditAgent(row)" @keydown.space.prevent="goEditAgent(row)">
                  {{ row.name || row.agent_name || '-' }}
                </span>
              </template>
            </el-table-column>
            <el-table-column :label="$t('auto.KBDrawer.cc6c35a3')" width="140" show-overflow-tooltip>
              <template #default="{ row }">
                {{ row.agent_code || row.code || '-' }}
              </template>
            </el-table-column>
            <el-table-column :label="$t('auto.KBDrawer.226b0912')" width="110" align="center">
              <template #default="{ row }">
                <el-tag size="small">{{ row.agent_type || row.type || '-' }}</el-tag>
              </template>
            </el-table-column>
            <el-table-column :label="$t('auto.KBDrawer.d129528d')" width="160">
              <template #default="{ row }">{{ formatTime(row.created_at || row.bind_at) }}</template>
            </el-table-column>
            <el-table-column :label="$t('auto.KBDrawer.2b6bc0f2')" width="100" align="center" fixed="right">
              <template #default="{ row }">
                <el-button link type="primary" size="small" @click="goEditAgent(row)">
                  查看
                </el-button>
              </template>
            </el-table-column>
          </el-table>
        </el-card>

        
        <el-card shadow="never" class="section">
          <template #header>
            <div class="section-header">
              <el-icon><List /></el-icon>
              <span>{{ statsLabel.item }}预览（前 20 条）</span>
            </div>
          </template>
          <el-empty
            v-if="!loadingItems && itemsPreview.length === 0"
            :description="$t('auto.KBDrawer.438afb7b')"
            :image-size="80"
          />
          <el-table v-else :data="itemsPreview" stripe size="small" border max-height="320">
            <el-table-column
              v-for="col in itemColumns"
              :key="col.prop"
              :prop="col.prop"
              :label="col.label"
              :min-width="col.minWidth || 120"
              :width="col.width"
              show-overflow-tooltip
            />
          </el-table>
        </el-card>
      </template>
      <el-empty v-else :description="$t('auto.KBDrawer.1a5036dd')" />
    </div>

    <template #footer>
      <el-button @click="$emit('update:modelValue', false)">{{ $t('auto.KBDrawer.b15d9127') }}</el-button>
      <el-button
        type="primary"
        @click="goEdit"
        v-if="kbData?.id"
      >
        <el-icon><Edit /></el-icon>编辑知识库
      </el-button>
    </template>
  </el-drawer>
</template>

<script setup>
import { ref, computed, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import {
  InfoFilled,
  DataLine,
  Connection,
  List,
  Refresh,
  Edit
} from '@element-plus/icons-vue'
import { getKB, getKBStats } from '@/api/knowledgeBase'
import { listByKB } from '@/api/agentKBBinding'

const props = defineProps({
  modelValue: { type: Boolean, default: false },
  kbData: { type: Object, default: null }
})

defineEmits(['update:modelValue', 'updated'])
const router = useRouter()

const { t } = useI18n()

const loading = ref(false)
const loadingAgents = ref(false)
const loadingItems = ref(false)
const agentBindings = ref([])
const itemsPreview = ref([])
const stats = ref({ agent_count: 0 })

const drawerTitle = computed(() => {
  if (!props.kbData) return t('auto.KBDrawer.cbdd5794')
  return `知识库详情：${props.kbData.name || props.kbData.kb_code || '#' + props.kbData.id}`
})

const statsLabel = computed(() => {
  const t = props.kbData?.kb_type
  if (t === 'rag') return { item: t('auto.KBDrawer.6bbc6401') }
  if (t === 'faq') return { item: 'FAQ 条目' }
  if (t === 'sop') return { item: 'SOP 模板' }
  return { item: t('auto.KBDrawer.ec39bb3b') }
})

const itemColumns = computed(() => {
  const t = props.kbData?.kb_type
  if (t === 'rag') {
    return [
      { prop: 'id', label: 'ID', width: 70 },
      { prop: 'title', label: t('auto.KBDrawer.bf13680e'), minWidth: 180 },
      { prop: 'file_type', label: t('auto.KBDrawer.226b0912'), width: 90 },
      { prop: 'chunk_count', label: t('auto.KBDrawer.f3cb638c'), width: 90 },
      { prop: 'created_at', label: t('auto.KBDrawer.eca37cb0'), width: 160 }
    ]
  }
  if (t === 'faq') {
    return [
      { prop: 'id', label: 'ID', width: 70 },
      { prop: 'question', label: t('auto.KBDrawer.5dc99f6e'), minWidth: 200 },
      { prop: 'category', label: t('auto.KBDrawer.d0771a42'), width: 120 },
      { prop: 'intent', label: t('auto.KBDrawer.764aca4e'), width: 120 },
      { prop: 'enabled', label: t('auto.KBDrawer.3fea7ca7'), width: 80 }
    ]
  }
  if (t === 'sop') {
    return [
      { prop: 'id', label: 'ID', width: 70 },
      { prop: 'name', label: t('auto.KBDrawer.a5d1c511'), minWidth: 180 },
      { prop: 'stage', label: t('auto.KBDrawer.d316330b'), width: 100 },
      { prop: 'intent', label: t('auto.KBDrawer.764aca4e'), width: 120 },
      { prop: 'enabled', label: t('auto.KBDrawer.3fea7ca7'), width: 80 }
    ]
  }
  return [
    { prop: 'id', label: 'ID', width: 70 },
    { prop: 'name', label: t('auto.KBDrawer.d7ec2d3f'), minWidth: 180 }
  ]
})

const getTypeLabel = (type) => {
  const map = { rag: 'RAG 文档', faq: 'FAQ', sop: 'SOP 模板' }
  return map[type] || type || '-'
}

const getTypeTagType = (type) => {
  const map = { rag: 'primary', faq: 'success', sop: 'warning' }
  return map[type] || 'info'
}

const formatTime = (t) => {
  if (!t) return '-'
  try {
    const d = new Date(t)
    if (isNaN(d.getTime())) return String(t)
    const pad = (n) => String(n).padStart(2, '0')
    return `${d.getFullYear()}-${pad(d.getMonth() + 1)}-${pad(d.getDate())} ${pad(d.getHours())}:${pad(d.getMinutes())}`
  } catch {
    return String(t)
  }
}

const loadAgentBindings = async () => {
  if (!props.kbData?.id) return
  loadingAgents.value = true
  try {
    const res = await listByKB(props.kbData.id).catch(() => null)
    const list = Array.isArray(res) ? res : res?.list || res?.items || []
    agentBindings.value = list
    stats.value.agent_count = list.length
  } catch (e) {
    agentBindings.value = []
  } finally {
    loadingAgents.value = false
  }
}

const loadStats = async () => {
  if (!props.kbData?.id) return
  try {
    const res = await getKBStats(props.kbData.id).catch(() => null)
    if (res) {
      stats.value = {
        agent_count: res.agent_count ?? stats.value.agent_count
      }
    }
  } catch { /* 忽略：清理/存储/恢复类 best-effort 操作 */ }
}

const loadItemsPreview = async () => {
  if (!props.kbData?.id) return
  loadingItems.value = true
  try {
    const t = props.kbData.kb_type;
    let url = ''
    if (t === 'rag') url = '/api/knowledge'
    else if (t === 'faq') url = '/api/faqs'
    else if (t === 'sop') url = '/api/sop-templates'
    if (!url) {
      itemsPreview.value = []
      return
    }
    const { http } = await import('@/utils/request')
    const res = await http
      .get(url, { kb_id: props.kbData.id, page: 1, page_size: 20 })
      .catch(() => null)
    const list = Array.isArray(res) ? res : res?.list || res?.items || []
    itemsPreview.value = list.slice(0, 20)
  } catch {
    itemsPreview.value = []
  } finally {
    loadingItems.value = false
  }
}

const goEdit = () => {
  if (!props.kbData?.id) return
  router.push({ name: 'KnowledgeBaseEdit', params: { id: props.kbData.id } })
}

const goEditAgent = (row) => {
  const id = row.agent_id || row.id
  if (!id) return
  router.push({ name: 'AIAgentEdit', params: { id } })
}

watch(
  () => [props.modelValue, props.kbData?.id],
  ([visible, id]) => {
    if (visible && id) {
      loading.value = true
      Promise.all([loadAgentBindings(), loadStats(), loadItemsPreview()])
        .catch(() => {})
        .finally(() => {
          loading.value = false
        })
    }
  },
  { immediate: true }
)
</script>

<style scoped>
.kb-drawer {
  padding: 0 4px;
}
.section {
  margin-bottom: 16px;
}
.section-header {
  display: flex;
  align-items: center;
  gap: 6px;
  font-weight: 600;
}
.mini-stat {
  text-align: center;
  padding: 12px;
  background: #f5f7fa;
  border-radius: 6px;
}
.mini-stat-label {
  color: #909399;
  font-size: 12px;
}
.mini-stat-value {
  font-size: 22px;
  font-weight: 600;
  margin-top: 4px;
  color: #303133;
}
.mini-stat-hint {
  margin-top: 2px;
  font-size: 11px;
  color: #c0c4cc;
}
.text-primary {
  color: #409eff;
}
.text-success {
  color: #67c23a;
}
.agent-link {
  color: #409eff;
  cursor: pointer;
}
.agent-link:hover {
  text-decoration: underline;
}
</style>
