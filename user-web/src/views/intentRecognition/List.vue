<template>
  <div class="intent-recognition-page">
    
    <el-card class="header-card" shadow="never">
      <div class="header-content">
        <div>
          <h2>{{ $t('销售意图识别') }}</h2>
          <p class="subtitle">{{ $t('auto.List.b5104421') }}</p>
        </div>
        <div class="header-actions">
          <el-tooltip
            ::content="$t('auto.List.e8cf410c')"
            placement="bottom"
            :show-after="200"
          >
            <el-switch
              v-model="intentEnabled"
              :loading="configLoading"
              inline-prompt
              :active-text="$t('auto.List.7854b52a')"
              :inactive-text="$t('auto.List.b15d9127')"
              style="margin-right: 12px"
              @change="onToggleIntent"
            />
          </el-tooltip>
          <el-button @click="openKeywordOverride" style="margin-right: 12px">
            <el-icon><EditPen /></el-icon>
            {{ $t('词表配置') }}
          </el-button>
          <el-button @click="refreshAll" :loading="loading">
            <el-icon><Refresh /></el-icon>
            {{ $t('刷新') }}
          </el-button>
        </div>
      </div>
      <div v-if="configMeta.updated_at" class="config-meta">
        <el-tag size="small" effect="plain" type="info">
          最近更新：{{ formatTime(configMeta.updated_at) }} · by {{ configMeta.updated_by || 'system' }}
        </el-tag>
      </div>
    </el-card>

    
    <el-card shadow="never" class="test-card">
      <template #header>
        <div class="card-header">
          <span><el-icon><Aim /></el-icon> {{ $t('意图识别测试') }}</span>
          <el-button link type="primary" @click="fillExample">{{ $t('填入示例') }}</el-button>
        </div>
      </template>
      <el-form :model="testForm" label-width="90px">
        <el-row :gutter="16">
          <el-col :span="12">
            <el-form-item :label="$t('auto.List.3f61dec5')">
              <el-input v-model="testForm.customer_id" :placeholder="$t('auto.List.c080949d')" clearable />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item :label="$t('auto.List.4981287f')">
              <el-select v-model="testForm.platform" :placeholder="$t('auto.List.32c3da83')" clearable style="width: 100%">
                <el-option :label="$t('auto.List.cfbf6f4c')" value="wechat" />
                <el-option :label="$t('auto.List.277a984f')" value="douyin" />
                <el-option :label="$t('auto.List.ae26bb04')" value="xiaohongshu" />
                <el-option :label="$t('auto.List.f32c0a2a')" value="kuaishou" />
                <el-option label="WhatsApp" value="whatsapp" />
                <el-option label="Telegram" value="telegram" />
              </el-select>
            </el-form-item>
          </el-col>
        </el-row>
        <el-form-item :label="$t('auto.List.50f198f0')">
          <el-input v-model="testForm.context" :placeholder="$t('auto.List.b0f349c4')" clearable />
        </el-form-item>
        <el-form-item :label="$t('auto.List.80ded089')">
          <el-input
            v-model="testForm.message"
            type="textarea"
            :rows="3"
            :placeholder="$t('auto.List.ac27ae19')"
          />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" :loading="recognizing" @click="runRecognize">
            <el-icon><Aim /></el-icon> 开始识别
          </el-button>
          <el-button type="success" :loading="batchLoading" @click="runBatchRecognize">{{ $t('auto.List.ab0467e1') }}</el-button>
        </el-form-item>
      </el-form>

      <template v-if="recognizeResult">
        <el-divider content-position="left">{{ $t('auto.List.fcf276af') }}</el-divider>
        <el-descriptions :column="2" border>
          <el-descriptions-item :label="$t('auto.List.766c070f')">
            <el-tag :type="getIntentTagType(recognizeResult.intent_type)">
              {{ recognizeResult.intent_name || intentNameMap[recognizeResult.intent_type] || recognizeResult.intent_type || '-' }}
            </el-tag>
          </el-descriptions-item>
          <el-descriptions-item :label="$t('auto.List.8f3599b6')">
            <el-progress
              :percentage="Math.round((recognizeResult.confidence || 0) * 100)"
              :color="getConfidenceColor(recognizeResult.confidence_level)"
              :stroke-width="14"
              :text-inside="true"
              style="width: 220px"
            />
          </el-descriptions-item>
          <el-descriptions-item :label="$t('auto.List.eebd3f23')">
            <el-tag :type="getConfidenceTagType(recognizeResult.confidence_level)" size="small" effect="plain">
              {{ recognizeResult.confidence_level || '-' }}
            </el-tag>
          </el-descriptions-item>
          <el-descriptions-item :label="$t('auto.List.6fcade01')">
            <el-tag :type="recognizeResult.method === 'llm' ? 'danger' : 'info'" size="small" effect="plain">
              {{ recognizeResult.method === 'llm' ? 'LLM 模型' : '规则匹配' }}
            </el-tag>
          </el-descriptions-item>
          <el-descriptions-item v-if="recognizeResult.platform" :label="$t('auto.List.4981287f')">
            {{ recognizeResult.platform }}
          </el-descriptions-item>
          <el-descriptions-item v-if="recognizeResult.latency_ms" :label="$t('auto.List.39f1374d')">
            {{ recognizeResult.latency_ms }} ms
          </el-descriptions-item>
        </el-descriptions>

        <template v-if="recognizeResult.entities && Object.keys(recognizeResult.entities).length">
          <el-divider content-position="left">{{ $t('auto.List.d3d9ad96') }}</el-divider>
          <div class="entity-list">
            <el-tag v-for="(val, key) in recognizeResult.entities" :key="key" effect="plain" class="entity-tag">
              {{ key }}: {{ val }}
            </el-tag>
          </div>
        </template>
      </template>

      <template v-if="batchResults.length">
        <el-divider content-position="left">批量识别结果（{{ batchResults.length }} 条）</el-divider>
        <el-table :data="batchResults" stripe size="small" max-height="320">
          <el-table-column :label="$t('auto.List.932ccc57')" min-width="200" show-overflow-tooltip>
            <template #default="{ row }">{{ row.message || row.raw_text }}</template>
          </el-table-column>
          <el-table-column :label="$t('auto.List.764aca4e')" width="120">
            <template #default="{ row }">
              <el-tag :type="getIntentTagType(row.intent_type)" size="small">
                {{ row.intent_name || intentNameMap[row.intent_type] || row.intent_type }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column :label="$t('auto.List.8f3599b6')" width="100">
            <template #default="{ row }">{{ Math.round((row.confidence || 0) * 100) }}%</template>
          </el-table-column>
          <el-table-column :label="$t('auto.List.ea340b9d')" width="80">
            <template #default="{ row }">
              <el-tag :type="getMethodTagType(row.method)" size="small" effect="plain">
                {{ getMethodLabel(row.method) }}
              </el-tag>
            </template>
          </el-table-column>
        </el-table>
      </template>
    </el-card>

    
    <el-row :gutter="16" class="stats-row" v-loading="statsLoading">
      <el-col :span="6">
        <el-card shadow="hover" class="stat-card">
          <div class="stat-value">{{ statsData.total || 0 }}</div>
          <div class="stat-label">{{ $t('auto.List.2616b1f3') }}</div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card shadow="hover" class="stat-card stat-purchase">
          <div class="stat-value">{{ statsData.purchase || 0 }}</div>
          <div class="stat-label">{{ $t('auto.List.ce092ef3') }}</div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card shadow="hover" class="stat-card stat-inquiry">
          <div class="stat-value">{{ statsData.price_inquiry || 0 }}</div>
          <div class="stat-label">{{ $t('auto.List.be446cfe') }}</div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card shadow="hover" class="stat-card stat-objection">
          <div class="stat-value">{{ objectionTotal }}</div>
          <div class="stat-label">{{ $t('auto.List.ec0ea0bd') }}</div>
        </el-card>
      </el-col>
    </el-row>

    <el-row :gutter="16">
      
      <el-col :span="10">
        <el-card shadow="never">
          <template #header>
            <div class="card-header">
              <span>{{ $t('auto.List.549a894b') }}</span>
              <el-button link type="primary" size="small" @click="loadStats">{{ $t('auto.List.694fc5ef') }}</el-button>
            </div>
          </template>
          <div class="intent-distribution" v-loading="statsLoading">
            <template v-if="distributionList.length">
              <div v-for="item in distributionList" :key="item.type" class="distribution-item">
                <span class="dist-name">
                  <el-tag :type="getIntentTagType(item.type)" size="small" effect="plain">
                    {{ intentNameMap[item.type] || item.type }}
                  </el-tag>
                </span>
                <div class="bar-track">
                  <div class="bar-fill" :style="{ width: distBarWidth(item.count) + '%', background: getIntentColor(item.type) }"></div>
                </div>
                <span class="dist-count">{{ item.count }}</span>
              </div>
            </template>
            <el-empty v-else :description="$t('auto.List.a3b437c9')" :image-size="80" />
          </div>
        </el-card>
      </el-col>

      
      <el-col :span="14">
        <el-card shadow="never">
          <template #header>
            <div class="card-header">
              <span>{{ $t('auto.List.2ca96e39') }}</span>
              <div class="header-controls">
                <el-select v-model="filterIntentType" :placeholder="$t('auto.List.c212e922')" clearable size="small" style="width: 160px" @change="onFilterChange">
                  <el-option v-for="d in intentDict" :key="d.type" :label="d.name" :value="d.type" />
                </el-select>
                <el-button size="small" type="primary" @click="loadRecent">{{ $t('auto.List.bee912d7') }}</el-button>
              </div>
            </div>
          </template>
          <el-table :data="recentList" v-loading="recentLoading" stripe size="small" style="width: 100%">
            <el-table-column :label="$t('auto.List.b87b7756')" min-width="180" show-overflow-tooltip>
              <template #default="{ row }">{{ row.message || row.raw_text }}</template>
            </el-table-column>
            <el-table-column :label="$t('auto.List.c6d069c9')" width="120">
              <template #default="{ row }">
                <el-tag :type="getIntentTagType(row.intent_type)" size="small">
                  {{ intentNameMap[row.intent_type] || row.intent_type }}
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column :label="$t('auto.List.8f3599b6')" width="130">
              <template #default="{ row }">
                <el-progress
                  :percentage="Math.round((row.confidence || 0) * 100)"
                  :color="getConfidenceColor(row.confidence_level)"
                  :stroke-width="12"
                  :text-inside="true"
                />
              </template>
            </el-table-column>
            <el-table-column :label="$t('auto.List.4981287f')" width="90">
              <template #default="{ row }">
                <el-tag v-if="row.platform" size="small" effect="plain">{{ getChannelLabel(row.platform) }}</el-tag>
                <span v-else>-</span>
              </template>
            </el-table-column>
            <el-table-column :label="$t('auto.List.19fcb9eb')" width="150">
              <template #default="{ row }">{{ formatTime(row.created_at) }}</template>
            </el-table-column>
            <template #empty>
              <el-empty :description="$t('auto.List.b9db8476')" :image-size="80" />
            </template>
          </el-table>
          <div class="pagination-container" v-if="recentPagination.total > 0">
            <el-pagination
              v-model:current-page="recentPagination.page"
              v-model:page-size="recentPagination.pageSize"
              :page-sizes="[10, 20, 50]"
              layout="total, prev, pager, next"
              :total="recentPagination.total"
              @current-change="loadRecent"
              @size-change="loadRecent"
            />
          </div>
        </el-card>
      </el-col>
    </el-row>

    
    <el-card shadow="never" class="test-card">
      <template #header>
        <div class="card-header">
          <span><el-icon><MagicStick /></el-icon>{{ $t('auto.List.902eda12') }}</span>
          <el-button link type="primary" @click="fillFineExample">{{ $t('填入示例') }}</el-button>
        </div>
      </template>
      <el-form :model="fineForm" label-width="90px" inline>
        <el-form-item :label="$t('auto.List.80ded089')">
          <el-input v-model="fineForm.message" type="textarea" :rows="2" :placeholder="$t('auto.List.c8be063c')" style="width: 460px" />
        </el-form-item>
        <el-form-item :label="$t('auto.List.3f61dec5')">
          <el-input v-model="fineForm.customer_id" :placeholder="$t('auto.List.c20cba89')" clearable style="width: 180px" />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" :loading="fineLoading" @click="runRecognizeFine">
            <el-icon><MagicStick /></el-icon> 精细识别
          </el-button>
        </el-form-item>
      </el-form>
      <template v-if="fineResult">
        <el-divider content-position="left">{{ $t('auto.List.20d4ea2c') }}</el-divider>
        <el-descriptions :column="2" border>
          <el-descriptions-item :label="$t('auto.List.d196c728')">
            <el-tag type="primary" effect="dark">{{ fineResult.major || '-' }}</el-tag>
          </el-descriptions-item>
          <el-descriptions-item :label="$t('auto.List.815e210d')">
            <el-tag effect="plain">{{ fineResult.minor || '-' }}</el-tag>
          </el-descriptions-item>
          <el-descriptions-item :label="$t('auto.List.8f3599b6')">
            <el-progress
              :percentage="Math.round((fineResult.confidence || 0) * 100)"
              :color="(fineResult.confidence || 0) >= 0.7 ? '#10B981' : '#F59E0B'"
              :stroke-width="14"
              style="width: 220px"
            />
          </el-descriptions-item>
          <el-descriptions-item :label="$t('auto.List.ea340b9d')">
            <el-tag :type="getFineMethodType(fineResult.method)" size="small">
              {{ getFineMethodLabel(fineResult.method) }}
            </el-tag>
          </el-descriptions-item>
          <el-descriptions-item v-if="fineResult.llm_model" :label="$t('auto.List.c82b0d96')">
            <el-tag size="small" effect="plain">{{ fineResult.llm_model }}</el-tag>
          </el-descriptions-item>
          <el-descriptions-item v-if="fineResult.latency_ms" :label="$t('auto.List.39f1374d')">
            {{ fineResult.latency_ms }} ms
          </el-descriptions-item>
          <el-descriptions-item v-if="fineResult.reasoning" :label="$t('auto.List.633fb3ff')" :span="2">
            {{ fineResult.reasoning }}
          </el-descriptions-item>
        </el-descriptions>
      </template>
    </el-card>

    
    <el-card shadow="never" class="dict-card">
      <template #header>
        <div class="card-header">
          <span>精细意图识别日志（最近 {{ fineLogs.length }} 条）</span>
          <div class="header-controls">
            <el-select v-model="filterFineMajor" :placeholder="$t('auto.List.30ad2af2')" clearable size="small" style="width: 160px" @change="loadFineLogs">
              <el-option v-for="m in fineMajorOptions" :key="m" :label="m" :value="m" />
            </el-select>
            <el-button size="small" type="primary" @click="loadFineLogs">{{ $t('auto.List.bee912d7') }}</el-button>
          </div>
        </div>
      </template>
      <el-table :data="fineLogs" v-loading="fineLogsLoading" stripe size="small">
        <el-table-column :label="$t('auto.List.b87b7756')" min-width="200" show-overflow-tooltip>
          <template #default="{ row }">{{ row.message || '-' }}</template>
        </el-table-column>
        <el-table-column :label="$t('auto.List.d196c728')" width="120">
          <template #default="{ row }">
            <el-tag type="primary" size="small" effect="plain">{{ row.intent_major || '-' }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column :label="$t('auto.List.815e210d')" width="160" show-overflow-tooltip>
          <template #default="{ row }">{{ row.intent_minor || '-' }}</template>
        </el-table-column>
        <el-table-column :label="$t('auto.List.8f3599b6')" width="100">
          <template #default="{ row }">{{ Math.round((row.confidence || 0) * 100) }}%</template>
        </el-table-column>
        <el-table-column :label="$t('auto.List.ea340b9d')" width="80">
          <template #default="{ row }">
            <el-tag :type="getFineMethodType(row.method)" size="small" effect="plain">
              {{ getFineMethodLabel(row.method) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column :label="$t('auto.List.39f1374d')" width="80">
          <template #default="{ row }">{{ row.latency_ms || 0 }}ms</template>
        </el-table-column>
        <el-table-column :label="$t('auto.List.19fcb9eb')" width="150">
          <template #default="{ row }">{{ formatTime(row.timestamp || row.created_at) }}</template>
        </el-table-column>
        <template #empty>
          <el-empty :description="$t('auto.List.04c99b5a')" :image-size="80" />
        </template>
      </el-table>
    </el-card>

    
    <el-card shadow="never" class="dict-card">
      <template #header>
        <div class="card-header">
          <span>意图词典（{{ intentDict.length }} 类）</span>
          <el-input v-model="dictKeyword" :placeholder="$t('auto.List.f12a4026')" clearable size="small" style="width: 220px" />
        </div>
      </template>
      <el-table :data="filteredDict" stripe size="small">
        <el-table-column :label="$t('auto.List.764aca4e')" width="130">
          <template #default="{ row }">
            <el-tag :type="getIntentTagType(row.type)" size="small">{{ row.name }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="type" :label="$t('auto.List.45235f6c')" width="180" />
        <el-table-column prop="description" :label="$t('auto.List.f411d0f1')" min-width="180" show-overflow-tooltip />
        <el-table-column :label="$t('auto.List.9699a50e')" min-width="260">
          <template #default="{ row }">
            <el-tag v-for="kw in (row.keywords || []).slice(0, 6)" :key="kw" size="small" effect="plain" style="margin: 2px">
              {{ kw }}
            </el-tag>
            <span v-if="(row.keywords || []).length > 6" class="more-tag">等 {{ row.keywords.length }} 个</span>
          </template>
        </el-table-column>
        <el-table-column :label="$t('auto.List.1a63ac23')" min-width="220" show-overflow-tooltip>
          <template #default="{ row }">
            {{ (row.examples || []).join(' / ') }}
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <!-- 词表配置对话框 -->
    <el-dialog v-model="keywordOverrideVisible" :title="$t('auto.List.24690feb')" width="720px">
      <el-alert
        :title="$t('auto.List.38a2e6dc')"
        type="info"
        :closable="false"
        show-icon
        style="margin-bottom: 16px"
      />
      <div v-if="keywordOverrideLoading" style="text-align: center; padding: 40px">
        <el-icon class="is-loading" :size="32"><Loading /></el-icon>
        <p style="margin-top: 12px; color: #909399">{{ $t('auto.List.26b5bd49') }}</p>
      </div>
      <div v-else>
        <div v-for="intent in keywordOverrideIntents" :key="intent.type" style="margin-bottom: 20px">
          <div style="display: flex; align-items: center; margin-bottom: 8px">
            <el-tag :type="getIntentTagType(intent.type)" size="small">
              {{ intent.name }}
            </el-tag>
            <span style="margin-left: 8px; font-size: 12px; color: #909399">({{ intent.type }})</span>
          </div>
          <div style="display: flex; flex-wrap: wrap; gap: 8px; margin-bottom: 8px">
            <el-tag
              v-for="(kw, idx) in (keywordOverrideDraft[intent.type] || [])"
              :key="idx"
              closable
              @close="removeKeyword(intent.type, idx)"
              size="small"
              effect="plain"
            >
              {{ kw }}
            </el-tag>
            <el-tag v-if="!keywordOverrideDraft[intent.type]?.length" type="info" size="small" effect="plain">
              暂无覆盖词
            </el-tag>
          </div>
          <div style="display: flex; gap: 8px">
            <el-input
              v-model="keywordInput[intent.type]"
              :placeholder="$t('auto.List.5b5c82bd')"
              size="small"
              style="flex: 1"
              @keyup.enter="addKeyword(intent.type)"
            />
            <el-button size="small" @click="addKeyword(intent.type)">{{ $t('auto.List.b58c7549') }}</el-button>
          </div>
        </div>
      </div>
      <template #footer>
        <el-button @click="keywordOverrideVisible = false">{{ $t('auto.List.625fb26b') }}</el-button>
        <el-button type="primary" :loading="keywordOverrideSaving" @click="saveKeywordOverride">{{ $t('auto.List.be5fbbe3') }}</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import i18n from '@/i18n'

import { ref, computed, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { ElMessage } from 'element-plus'
import { Aim, MagicStick, Refresh, Loading, EditPen } from '@element-plus/icons-vue'
import { intentApi } from '@/api/intentRecognition.js'
import { getChannelLabel } from '@/constants/channel'
const RECOGNIZE_METHOD = { llm: 'LLM', rule: t('auto.List.b0fae043'), keyword: t('auto.List.9699a50e'), bert: 'BERT', hybrid: t('auto.List.41ce131a') };
const RECOGNIZE_METHOD_TAG = { llm: 'danger', rule: 'info', keyword: 'success', bert: 'warning', hybrid: 'primary' }
const getMethodLabel = (m) => RECOGNIZE_METHOD[m] || m || '-'
// 兜底必须是 el-tag 的合法 type：空串会触发
// "Invalid prop: validation failed for prop type" 告警并回落到 default(primary) 显示成蓝色。
// 与本文件 getIntentTagType(:666) / getFineMethodType(:517) 的 'info' 兜底一致。
const getMethodTagType = (m) => RECOGNIZE_METHOD_TAG[m] || 'info'

const { t } = useI18n()

const statsLoading = ref(false);
const statsData = ref({})

const recentLoading = ref(false);
const recentList = ref([])
const filterIntentType = ref('')
const recentPagination = ref({ page: 1, pageSize: 20, total: 0 })

const intentDict = ref([]);
const dictKeyword = ref('')

const recognizing = ref(false);
const testForm = ref({ customer_id: '', platform: '', context: '', message: '' })
const recognizeResult = ref(null)

const batchLoading = ref(false);
const batchResults = ref([])

const fineLoading = ref(false);
const fineForm = ref({ message: '', customer_id: '' })
const fineResult = ref(null)
const fineLogsLoading = ref(false)
const fineLogs = ref([])
const filterFineMajor = ref('')
const fineMajorOptions = [
  'consult', 'price_inquiry', 'objection', 'after_sale',
  'complaint', 'churn', 'intent_buy', 'ask_product'
];
const FINE_METHOD_LABEL = { rule: t('auto.List.b0fae043'), llm: 'LLM', hybrid: t('auto.List.2e5edb3c') }
const FINE_METHOD_TAG = { rule: 'success', llm: 'danger', hybrid: 'warning' }
const getFineMethodLabel = (m) => FINE_METHOD_LABEL[m] || m || '-'
const getFineMethodType = (m) => FINE_METHOD_TAG[m] || 'info'

const loading = ref(false)

const configLoading = ref(false);
const intentEnabled = ref(true)
const configMeta = ref({ updated_at: '', updated_by: '' })

const loadIntentConfig = async () => {
  configLoading.value = true
  try {
    const res = await intentApi.getConfig()
    if (res && typeof res === 'object') {
      intentEnabled.value = res.enabled !== false
      configMeta.value = {
        updated_at: res.updated_at || '',
        updated_by: res.updated_by || ''
      }
    }
  } catch (e) {
    intentEnabled.value = true;
    ElMessage.warning(t('auto.List.2c5ca90b') + (e.message || t('auto.List.974e7484')))
  } finally {
    configLoading.value = false
  }
};

const onToggleIntent = async (val) => {
  configLoading.value = true
  try {
    const res = await intentApi.updateConfig({ enabled: val })
    if (res && typeof res === 'object') {
      configMeta.value = {
        updated_at: res.updated_at || '',
        updated_by: res.updated_by || ''
      }
    }
    ElMessage.success(`意图识别已${val ? t('auto.List.7854b52a') : t('auto.List.b15d9127')}，业务链路立即生效`)
  } catch (e) {
    intentEnabled.value = !val;
    ElMessage.error(t('auto.List.d0759f5f') + (e.message || t('auto.List.974e7484')))
  } finally {
    configLoading.value = false
  }
};

// 词表配置对话框
const keywordOverrideVisible = ref(false)
const keywordOverrideLoading = ref(false)
const keywordOverrideSaving = ref(false)
const keywordOverrideDraft = ref({})
const keywordInput = ref({})

const keywordOverrideIntents = computed(() => {
  return intentDict.value.map(d => ({
    type: d.type,
    name: d.name
  }))
})

const openKeywordOverride = async () => {
  keywordOverrideVisible.value = true
  keywordOverrideLoading.value = true
  keywordOverrideDraft.value = {}
  keywordInput.value = {}
  
  try {
    const res = await intentApi.getKeywordOverride()
    if (res && typeof res === 'object') {
      keywordOverrideDraft.value = JSON.parse(JSON.stringify(res))
    }
  } catch (e) {
    ElMessage.error(t('auto.List.4c832b12') + (e.message || t('auto.List.974e7484')))
  } finally {
    keywordOverrideLoading.value = false
  }
}

const addKeyword = (type) => {
  const input = keywordInput.value[type]
  if (!input || !input.trim()) return
  
  if (!keywordOverrideDraft.value[type]) {
    keywordOverrideDraft.value[type] = []
  }
  
  const kw = input.trim()
  if (!keywordOverrideDraft.value[type].includes(kw)) {
    keywordOverrideDraft.value[type].push(kw)
  }
  keywordInput.value[type] = ''
}

const removeKeyword = (type, idx) => {
  if (keywordOverrideDraft.value[type]) {
    keywordOverrideDraft.value[type].splice(idx, 1)
  }
}

const saveKeywordOverride = async () => {
  keywordOverrideSaving.value = true
  try {
    await intentApi.updateKeywordOverride({ override: keywordOverrideDraft.value })
    ElMessage.success(t('auto.List.4e023cee'))
    keywordOverrideVisible.value = false
  } catch (e) {
    ElMessage.error(t('auto.List.036fe2d5') + (e.message || t('auto.List.974e7484')))
  } finally {
    keywordOverrideSaving.value = false
  }
}

const intentNameMap = computed(() => {
  const map = {}
  intentDict.value.forEach(d => { map[d.type] = d.name })
  return map
});

const objectionTotal = computed(() => {
  const s = statsData.value
  return (s.objection_price || 0) + (s.objection_need || 0) + (s.objection_trust || 0) +
    (s.objection_competitor || 0) + (s.objection_timing || 0)
})

const distributionList = computed(() => {
  const s = statsData.value
  if (Array.isArray(s.distribution)) {
    return s.distribution
  }
  return Object.entries(s)
    .filter(([k]) => !['total', 'distribution'].includes(k))
    .map(([type, count]) => ({ type, count }))
    .sort((a, b) => b.count - a.count)
})

const filteredDict = computed(() => {
  if (!dictKeyword.value) return intentDict.value
  const kw = dictKeyword.value.toLowerCase()
  return intentDict.value.filter(d =>
    d.name.includes(kw) || d.type.includes(kw) ||
    (d.keywords || []).some(k => k.includes(kw))
  )
})

const distBarWidth = (count) => {
  const max = Math.max(...distributionList.value.map(d => d.count), 1)
  return Math.round((count / max) * 100)
}

const getIntentTagType = (type) => {
  if (!type) return 'info'
  if (type === 'purchase') return 'success'
  if (type === 'price_inquiry' || type === 'ask_product') return 'primary'
  if (type.startsWith('objection')) return 'warning'
  if (type === 'churn' || type === 'complaint') return 'danger'
  if (type === 'unknown') return 'info'
  return 'info'
}

const getIntentColor = (type) => {
  const map = {
    purchase: '#10B981',
    price_inquiry: '#4F46E5',
    ask_product: '#4F46E5',
    ask_service: '#909399',
    after_sale: '#F59E0B',
    objection_price: '#F59E0B',
    objection_need: '#F59E0B',
    objection_trust: '#F59E0B',
    objection_competitor: '#F59E0B',
    objection_timing: '#F59E0B',
    churn: '#EF4444',
    complaint: '#EF4444',
    social: '#909399',
    greeting: '#909399',
    unknown: '#c0c4cc'
  }
  return map[type] || '#909399'
}

const getConfidenceColor = (level) => {
  if (level === 'high') return '#10B981'
  if (level === 'medium') return '#F59E0B'
  return '#EF4444'
}

// 置信度对应的 el-tag 语义色。必须与上面的 getConfidenceColor 分开：
// el-tag 的 type 只能取 primary/success/info/warning/danger，传 hex 颜色（实测返回 '#EF4444'）
// 会触发 "Invalid prop: validation failed for prop type" 告警，标签色还会错误回落到 primary 蓝色。
// 与本文件 getIntentTagType(:666) 的「颜色函数管图表、*TagType 管标签」配对写法一致。
const getConfidenceTagType = (level) => {
  if (level === 'high') return 'success'
  if (level === 'medium') return 'warning'
  return 'danger'
}

const formatTime = (t) => {
  if (!t) return '-'
  const d = new Date(t)
  if (isNaN(d.getTime())) return '-'
  const pad = (n) => String(n).padStart(2, '0')
  return `${d.getFullYear()}-${pad(d.getMonth() + 1)}-${pad(d.getDate())} ${pad(d.getHours())}:${pad(d.getMinutes())}`
}

const loadStats = async () => {
  statsLoading.value = true
  try {
    const res = await intentApi.getStats({
      page: 1,
      page_size: 100
    })
    if (res && typeof res === 'object') {
      statsData.value = res
    } else {
      statsData.value = {}
    }
  } catch (e) {
    ElMessage.error(t('auto.List.4f0b6e2e') + (e.message || t('auto.List.974e7484')))
  } finally {
    statsLoading.value = false
  }
}

const onFilterChange = () => {
  recentPagination.value.page = 1
  loadRecent()
}

const loadRecent = async () => {
  recentLoading.value = true
  try {
    const params = {
      page: recentPagination.value.page,
      page_size: recentPagination.value.pageSize
    }
    if (filterIntentType.value) params.intent_type = filterIntentType.value
    const res = await intentApi.getRecent(params)
    recentList.value = res?.list || (Array.isArray(res) ? res : []);
    recentPagination.value.total = res?.total || recentList.value.length
  } catch (e) {
    ElMessage.error(t('auto.List.b3272552') + (e.message || t('auto.List.974e7484')))
  } finally {
    recentLoading.value = false
  }
}

const loadDict = async () => {
  try {
    const res = await intentApi.getDict()
    intentDict.value = Array.isArray(res) ? res : (res?.list || [])
  } catch (e) {
    ElMessage.error(t('auto.List.a48487c4') + (e.message || t('auto.List.974e7484')))
  }
}

const refreshAll = async () => {
  loading.value = true
  try {
    await Promise.all([loadStats(), loadRecent(), loadDict()])
  } finally {
    loading.value = false
  }
}

const fillExample = () => {
  const examples = [
    t('auto.List.ce57caf7'),
    t('auto.List.ea54a824'),
    t('auto.List.16480225'),
    t('auto.List.ed2ed676'),
    t('auto.List.c1a0cf45')
  ]
  testForm.value.message = examples[Math.floor(Math.random() * examples.length)]
}

const runRecognize = async () => {
  if (!testForm.value.message || !testForm.value.message.trim()) {
    ElMessage.warning(i18n.global.t(t('auto.List.f6acdcf7')))
    return
  }
  recognizing.value = true
  batchResults.value = []
  try {
    const data = {
      message: testForm.value.message
    }
    if (testForm.value.context) data.context = testForm.value.context
    if (testForm.value.customer_id) data.customer_id = testForm.value.customer_id
    if (testForm.value.platform) data.platform = testForm.value.platform
    const res = await intentApi.recognize(data)
    recognizeResult.value = res
    ElMessage.success(i18n.global.t(t('auto.List.483f84a1')))
    loadRecent()
  } catch (e) {
    ElMessage.error(t('auto.List.a3ebed4c') + (e.message || t('auto.List.974e7484')))
  } finally {
    recognizing.value = false
  }
}

const runBatchRecognize = async () => {
  const lines = (testForm.value.message || '').split('\n').map(s => s.trim()).filter(Boolean);
  if (lines.length === 0) {
    ElMessage.warning(i18n.global.t(t('auto.List.41c20d3b')))
    return
  }
  batchLoading.value = true
  recognizeResult.value = null
  batchResults.value = []
  try {
    const res = await intentApi.batchRecognize({ messages: lines })
    const list = Array.isArray(res) ? res : (res?.list || [])
    batchResults.value = list.map((r, i) => ({ ...r, message: r.message || lines[i] }))
    ElMessage.success(`批量识别完成，共 ${batchResults.value.length} 条`)
    loadRecent()
  } catch (e) {
    ElMessage.error(t('auto.List.09d30f6a') + (e.message || t('auto.List.974e7484')))
  } finally {
    batchLoading.value = false
  }
}

const fillFineExample = () => {
  const examples = [
    t('auto.List.d9158db5'),
    t('auto.List.77041251'),
    t('auto.List.191b56b0'),
    t('auto.List.75f25947'),
    t('auto.List.5ed1ee42'),
    t('auto.List.e4dd3704')
  ]
  fineForm.value.message = examples[Math.floor(Math.random() * examples.length)]
};

const runRecognizeFine = async () => {
  if (!fineForm.value.message || !fineForm.value.message.trim()) {
    ElMessage.warning(i18n.global.t(t('auto.List.f6acdcf7')))
    return
  }
  fineLoading.value = true
  fineResult.value = null
  try {
    const data = { message: fineForm.value.message }
    if (fineForm.value.customer_id) data.customer_id = fineForm.value.customer_id
    const res = await intentApi.recognizeFine(data)
    fineResult.value = res
    ElMessage.success(i18n.global.t(t('auto.List.489a8195')))
    loadFineLogs()
  } catch (e) {
    ElMessage.error(t('auto.List.3569c4dd') + (e.message || t('auto.List.974e7484')))
  } finally {
    fineLoading.value = false
  }
}

const loadFineLogs = async () => {
  fineLogsLoading.value = true
  try {
    const params = { limit: 50 }
    if (filterFineMajor.value) params.major = filterFineMajor.value
    const res = await intentApi.getLogs(params)
    fineLogs.value = Array.isArray(res) ? res : (res?.list || [])
  } catch (e) {
    ElMessage.error(t('auto.List.aa7ea496') + (e.message || t('auto.List.974e7484')))
    fineLogs.value = []
  } finally {
    fineLogsLoading.value = false
  }
}

onMounted(() => {
  loadIntentConfig()
  refreshAll()
  loadFineLogs()
})
</script>

<style scoped lang="scss">
.intent-recognition-page { padding: 20px; }

.header-card {
  margin-bottom: 16px;
  :deep(.el-card__body) { padding: 16px 20px; }
  .header-content {
    display: flex;
    justify-content: space-between;
    align-items: center;
    h2 { margin: 0 0 6px 0; font-size: 20px; }
    .subtitle { color: #909399; margin: 0; font-size: 13px; }
  }
  .header-actions { display: flex; gap: 8px; }
}
.config-meta { margin-top: 12px; padding-top: 12px; border-top: 1px dashed #ebeef5; }

.test-card { margin-bottom: 16px; }

.stats-row {
  margin-bottom: 16px;
  .stat-card {
    text-align: center;
    .stat-value { font-size: 28px; font-weight: bold; color: #303133; }
    .stat-label { color: #909399; font-size: 13px; margin-top: 6px; }
    &.stat-purchase .stat-value { color: #10B981; }
    &.stat-inquiry .stat-value { color: #4F46E5; }
    &.stat-objection .stat-value { color: #F59E0B; }
  }
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  .header-controls { display: flex; gap: 8px; }
}

.intent-distribution {
  min-height: 200px;
  .distribution-item {
    display: flex;
    align-items: center;
    margin-bottom: 12px;
    .dist-name { width: 110px; flex-shrink: 0; }
    .bar-track {
      flex: 1;
      height: 14px;
      background: #f0f2f5;
      border-radius: 7px;
      overflow: hidden;
      margin: 0 10px;
      .bar-fill { height: 100%; border-radius: 7px; transition: width .3s; }
    }
    .dist-count { width: 40px; text-align: right; font-weight: bold; color: #606266; }
  }
}

.dict-card { margin-top: 16px; }
.more-tag { color: #909399; font-size: 12px; margin-left: 4px; }

.entity-list {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
  .entity-tag { margin: 2px; }
}

.pagination-container {
  margin-top: 16px;
  display: flex;
  justify-content: flex-end;
}
</style>
