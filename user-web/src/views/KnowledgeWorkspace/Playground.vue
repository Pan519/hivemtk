<template>
  <div class="playground-page">
    <el-alert
      ::title="$t('auto.Playground.05a06f61')"
      type="info"
      :closable="false"
      show-icon
    >
      <template #default>
        商户调试用:在生产前验证 topK、相似度阈值、过滤条件对检索质量的影响。命中分段可一键提交「相关/不相关」反馈,用于持续调优。
      </template>
    </el-alert>

    <el-row :gutter="16" style="margin-top: 16px">
      
      <el-col :span="8">
        <el-card>
          <template #header>
            <span>{{ $t('auto.Playground.78d2eb3d') }}</span>
          </template>
          <el-form :model="form" label-width="100px">
            <el-form-item :label="$t('auto.Playground.a015434e')" required>
              <el-select v-model="form.product_id" :placeholder="$t('auto.Playground.cbf1cfff')" style="width: 100%">
                <el-option v-for="p in productList" :key="p.id" :label="p.name" :value="p.id" />
              </el-select>
            </el-form-item>
            <el-form-item :label="$t('auto.Playground.c8ddfe05')" required>
              <el-input
                v-model="form.query"
                type="textarea"
                :rows="4"
                :placeholder="$t('auto.Playground.71393e9c')"
              />
            </el-form-item>
            <el-form-item label="Top K">
              <el-slider v-model="form.top_k" :min="1" :max="50" show-stops :marks="{ 5: '5', 10: '10', 20: '20', 50: '50' }" />
            </el-form-item>
            <el-form-item :label="$t('auto.Playground.2980caa0')">
              <el-slider v-model="form.similarity_threshold" :min="0" :max="1" :step="0.05" :format-tooltip="(v) => v.toFixed(2)" />
              <div class="threshold-hint">仅返回相似度 ≥ {{ form.similarity_threshold.toFixed(2) }} 的分段</div>
            </el-form-item>
            <el-form-item :label="$t('auto.Playground.4716d8d4')">
              <el-input v-model="form.filter_category" :placeholder="$t('auto.Playground.cbec3ea5')" clearable />
            </el-form-item>
            <el-form-item :label="$t('auto.Playground.5276f72e')">
              <el-select v-model="form.tags" multiple filterable allow-create :placeholder="$t('auto.Playground.26ff573c')" style="width: 100%">
                <el-option v-for="t in tagOptions" :key="t" :label="t" :value="t" />
              </el-select>
            </el-form-item>
            <el-form-item>
              <el-button type="primary" :icon="Search" :loading="searching" :disabled="!canSearch" @click="handleSearch" style="width: 100%">
                开始检索
              </el-button>
            </el-form-item>
            <el-form-item>
              <el-checkbox v-model="form.use_three_tier">{{ $t('auto.Playground.31ce2f52') }}</el-checkbox>
              <el-checkbox v-model="form.use_rerank">{{ $t('auto.Playground.2eeeabdf') }}</el-checkbox>
            </el-form-item>
          </el-form>
        </el-card>

        
        <el-card style="margin-top: 16px">
          <template #header>
            <span>{{ $t('auto.Playground.c4e86d88') }}</span>
          </template>
          <div class="template-list">
            <el-tag
              v-for="(q, idx) in templates"
              :key="idx"
              class="template-tag"
              @click="form.query = q"
            >{{ q }}</el-tag>
          </div>
        </el-card>
      </el-col>

      
      <el-col :span="16">
        <el-card v-loading="searching">
          <template #header>
            <div class="card-header">
              <span>{{ $t('auto.Playground.abb88af2') }}</span>
              <div v-if="result" class="result-stats">
                <el-tag size="small" type="info">命中 {{ result.total }} 条</el-tag>
                <!-- el-tag 的 type 合法值只有 primary/success/info/warning/danger。
                     传空串会触发 "Invalid prop: validation failed for prop type" 告警，
                     并回落到 default(primary) 显示成蓝色，与「未命中缓存」的中性语义不符。 -->
                <el-tag size="small" :type="result.from_cache ? 'success' : 'info'">
                  {{ result.from_cache ? '命中缓存' : '未命中缓存' }}
                </el-tag>
                <el-tag size="small" type="warning">耗时 {{ result.latency_ms }} ms</el-tag>
                <el-tag size="small">来源 {{ result.source }}</el-tag>
              </div>
            </div>
          </template>

          
          <el-row v-if="result" :gutter="12" class="metric-row">
            <el-col :span="6">
              <el-statistic :title="$t('auto.Playground.da21e885')" :value="result.max_score" :precision="3" />
            </el-col>
            <el-col :span="6">
              <el-statistic :title="$t('auto.Playground.e31d509e')" :value="result.min_score" :precision="3" />
            </el-col>
            <el-col :span="6">
              <el-statistic :title="$t('auto.Playground.3569877e')" :value="result.avg_score" :precision="3" />
            </el-col>
            <el-col :span="6">
              <el-statistic :title="$t('auto.Playground.1b68d85a')" :value="result.total" />
            </el-col>
          </el-row>

          <el-empty v-if="!result" :description="$t('auto.Playground.e297d57f')" />
          <div v-else-if="result.chunks.length === 0" class="no-hit">{{ $t('auto.Playground.5d4d6d01') }}</div>
          <div v-else class="chunks-list">
            <div v-for="(c, idx) in result.chunks" :key="c.chunk_id" class="chunk-card">
              <div class="chunk-header">
                <div class="chunk-meta">
                  <el-tag size="small" :type="getScoreType(c.score)">#{{ idx + 1 }} {{ c.score.toFixed(3) }}</el-tag>
                  <el-tag v-if="c.from_cache" size="small" type="success">{{ $t('auto.Playground.e80c310e') }}</el-tag>
                  <el-tag size="small" type="info">{{ c.source || 'L2' }}</el-tag>
                  <el-text size="small" type="info">chunk_id: {{ c.chunk_id }} | doc_id: {{ c.document_id }}</el-text>
                </div>
                <div class="chunk-actions">
                  <el-button-group>
                    <el-button size="small" :type="feedbacks[c.chunk_id] === 1 ? 'success' : ''" @click="submitFeedback(c, 1)">
                      <el-icon><CircleCheck /></el-icon> 相关
                    </el-button>
                    <el-button size="small" :type="feedbacks[c.chunk_id] === 0 ? 'info' : ''" @click="submitFeedback(c, 0)">
                      <el-icon><RemoveFilled /></el-icon> 一般
                    </el-button>
                    <el-button size="small" :type="feedbacks[c.chunk_id] === -1 ? 'danger' : ''" @click="submitFeedback(c, -1)">
                      <el-icon><CircleClose /></el-icon> 不相关
                    </el-button>
                  </el-button-group>
                  <el-button size="small" link type="primary" @click="goChunkEdit(c.document_id, c.chunk_id)">{{ $t('auto.Playground.87d4fe2d') }}</el-button>
                </div>
              </div>
              <div v-if="c.title" class="chunk-title">{{ c.title }}</div>
              <div class="chunk-content">{{ c.content }}</div>
              <el-progress
                v-if="c.score"
                :percentage="(c.score * 100).toFixed(1)"
                :stroke-width="4"
                :show-text="false"
                :color="getScoreColor(c.score)"
                style="margin-top: 6px"
              />
            </div>
          </div>
        </el-card>

        
        <el-card v-if="result && result.debug_info" style="margin-top: 16px">
          <template #header>
            <span>{{ $t('auto.Playground.a9cfab55') }}</span>
          </template>
          <el-descriptions :column="2" border size="small">
            <el-descriptions-item :label="$t('auto.Playground.2125533e')">{{ result.debug_info.product_id }}</el-descriptions-item>
            <el-descriptions-item label="Top K">{{ result.debug_info.top_k }}</el-descriptions-item>
            <el-descriptions-item :label="$t('auto.Playground.e0de1c9a')">{{ result.debug_info.threshold }}</el-descriptions-item>
            <el-descriptions-item :label="$t('auto.Playground.dd09412e')">{{ result.debug_info.use_three_tier ? '开启' : '关闭' }}</el-descriptions-item>
            <el-descriptions-item :label="$t('auto.Playground.f982377d')">{{ result.debug_info.use_rerank ? '开启' : '关闭' }}</el-descriptions-item>
            <el-descriptions-item :label="$t('auto.Playground.bee912d7')">{{ result.query }}</el-descriptions-item>
          </el-descriptions>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script setup>
import i18n from '@/i18n'

import { ref, reactive, computed, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { Search, CircleCheck, CircleClose, RemoveFilled } from '@element-plus/icons-vue'
import { knowledgeMerchantAPI } from '@/api/knowledgeMerchant'
import { ragProductConfigAPI } from '@/api/ragProductConfig'

const router = useRouter()

const { t } = useI18n()

const searching = ref(false)
const productList = ref([])
const result = ref(null)
const tagOptions = ref(['FAQ', t('auto.Playground.a015434e'), t('auto.Playground.59bd6897'), t('auto.Playground.b7824d5c'), t('auto.Playground.61a3ec66')])
const templates = [
  t('auto.Playground.23daa6aa'),
  t('auto.Playground.0591b740'),
  t('auto.Playground.50a572de'),
  t('auto.Playground.4fa5a2bf'),
  t('auto.Playground.e9868e22')
]

const feedbacks = reactive({})

const form = reactive({
  product_id: '',
  query: '',
  top_k: 5,
  similarity_threshold: 0,
  filter_category: '',
  tags: [],
  use_three_tier: false,
  use_rerank: false
})

const canSearch = computed(() => form.product_id && form.query.trim())

const loadProducts = async () => {
  try {
    const res = await ragProductConfigAPI.listProducts()
    if (Array.isArray(res)) productList.value = res
    else if (res?.items) productList.value = res.items
  } catch (e) {
    console.error(t('auto.Playground.4c46f3b1'), e)
  }
}

const handleSearch = async () => {
  if (!canSearch.value) {
    ElMessage.warning(i18n.global.t(t('auto.Playground.8bae35a9')))
    return
  }
  searching.value = true
  try {
    const res = await knowledgeMerchantAPI.playground({
      product_id: form.product_id,
      query: form.query,
      top_k: form.top_k,
      similarity_threshold: form.similarity_threshold,
      filter_category: form.filter_category,
      tags: form.tags,
      use_three_tier: form.use_three_tier,
      use_rerank: form.use_rerank
    })
    if (res) {
      result.value = res
      ElMessage.success(`检索完成,共 ${res.total} 条结果`)
    }
  } catch (e) {
    ElMessage.error(t('auto.Playground.203e2876') + (e.message || ''))
  } finally {
    searching.value = false
  }
}

const submitFeedback = async (chunk, rating) => {
  try {
    await knowledgeMerchantAPI.submitFeedback({
      product_id: form.product_id,
      query: form.query,
      document_id: chunk.document_id,
      chunk_id: chunk.chunk_id,
      rating: rating
    })
    feedbacks[chunk.chunk_id] = rating
    const labels = { 1: t('auto.Playground.ec0bce08'), 0: t('auto.Playground.2ab01e41'), '-1': t('auto.Playground.3541012e') }
    ElMessage.success(`已记录反馈: ${labels[rating] || rating}`)
  } catch (e) {
    ElMessage.error(t('auto.Playground.b6c39803') + (e.message || ''))
  }
}

const goChunkEdit = (documentId, chunkId) => {
  router.push({
    path: '/knowledge/management',
    query: { document_id: documentId, chunk_id: chunkId, tab: 'chunks' }
  })
}

const getScoreType = (score) => {
  if (score >= 0.8) return 'success'
  if (score >= 0.5) return 'warning'
  return 'danger'
}

const getScoreColor = (score) => {
  if (score >= 0.8) return '#10B981'
  if (score >= 0.5) return '#F59E0B'
  return '#EF4444'
}

onMounted(() => {
  loadProducts()
})
</script>

<style scoped lang="scss">
.playground-page {
  padding: 0;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.result-stats {
  display: flex;
  gap: 6px;
  flex-wrap: wrap;
}

.threshold-hint {
  font-size: 12px;
  color: #909399;
  margin-top: -10px;
}

.metric-row {
  margin-bottom: 12px;
  padding: 8px 0;
  border-bottom: 1px dashed #ebeef5;
}

.no-hit {
  text-align: center;
  color: #909399;
  padding: 40px 0;
}

.chunks-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
  margin-top: 12px;
}

.chunk-card {
  background: #fafafa;
  border: 1px solid #ebeef5;
  border-radius: 6px;
  padding: 12px;
  transition: all 0.2s;
  &:hover {
    border-color: #4F46E5;
    background: #f0f7ff;
  }
}

.chunk-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 8px;
  gap: 8px;
  flex-wrap: wrap;
}

.chunk-meta {
  display: flex;
  gap: 6px;
  align-items: center;
  flex-wrap: wrap;
}

.chunk-actions {
  display: flex;
  gap: 8px;
  align-items: center;
}

.chunk-title {
  font-size: 14px;
  font-weight: 600;
  color: #303133;
  margin-bottom: 4px;
}

.chunk-content {
  font-size: 13px;
  line-height: 1.6;
  color: #606266;
  white-space: pre-wrap;
  word-break: break-word;
}

.template-list {
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.template-tag {
  cursor: pointer;
  margin-right: 0;
}
</style>
