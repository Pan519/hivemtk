<template>
  <div class="intent-recognition-page">
    
    <el-card class="header-card" shadow="never">
      <div class="header-content">
        <div>
          <h2>{{ $t('销售意图识别') }}</h2>
          <p class="subtitle">{{ $t('auto.List.46.b51044') }}</p>
        </div>
        <div class="header-actions">
          <el-tooltip
            ::content="$t('auto.List.1.e8cf41')"
            placement="bottom"
            :show-after="200"
          >
            <el-switch
              v-model="intentEnabled"
              :loading="configLoading"
              inline-prompt
              :active-text="$t('auto.List.2.7854b5')"
              :inactive-text="$t('auto.List.3.b15d91')"
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
            <el-form-item :label="$t('auto.List.4.3f61de')">
              <el-input v-model="testForm.customer_id" :placeholder="$t('auto.List.5.c08094')" clearable />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item :label="$t('auto.List.6.498128')">
              <el-select v-model="testForm.platform" :placeholder="$t('auto.List.7.32c3da')" clearable style="width: 100%">
                <el-option :label="$t('auto.List.8.cfbf6f')" value="wechat" />
                <el-option :label="$t('auto.List.9.277a98')" value="douyin" />
                <el-option :label="$t('auto.List.10.ae26bb')" value="xiaohongshu" />
                <el-option :label="$t('auto.List.11.f32c0a')" value="kuaishou" />
                <el-option label="WhatsApp" value="whatsapp" />
                <el-option label="Telegram" value="telegram" />
              </el-select>
            </el-form-item>
          </el-col>
        </el-row>
        <el-form-item :label="$t('auto.List.12.50f198')">
          <el-input v-model="testForm.context" :placeholder="$t('auto.List.13.b0f349')" clearable />
        </el-form-item>
        <el-form-item :label="$t('auto.List.14.80ded0')">
          <el-input
            v-model="testForm.message"
            type="textarea"
            :rows="3"
            :placeholder="$t('auto.List.15.ac27ae')"
          />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" :loading="recognizing" @click="runRecognize">
            <el-icon><Aim /></el-icon> 开始识别
          </el-button>
          <el-button type="success" :loading="batchLoading" @click="runBatchRecognize">{{ $t('auto.List.47.ab0467') }}</el-button>
        </el-form-item>
      </el-form>

      <template v-if="recognizeResult">
        <el-divider content-position="left">{{ $t('auto.List.48.fcf276') }}</el-divider>
        <el-descriptions :column="2" border>
          <el-descriptions-item :label="$t('auto.List.16.766c07')">
            <el-tag :type="getIntentTagType(recognizeResult.intent_type)">
              {{ recognizeResult.intent_name || intentNameMap[recognizeResult.intent_type] || recognizeResult.intent_type || '-' }}
            </el-tag>
          </el-descriptions-item>
          <el-descriptions-item :label="$t('auto.List.17.8f3599')">
            <el-progress
              :percentage="Math.round((recognizeResult.confidence || 0) * 100)"
              :color="getConfidenceColor(recognizeResult.confidence_level)"
              :stroke-width="14"
              :text-inside="true"
              style="width: 220px"
            />
          </el-descriptions-item>
          <el-descriptions-item :label="$t('auto.List.18.eebd3f')">
            <el-tag :type="getConfidenceTagType(recognizeResult.confidence_level)" size="small" effect="plain">
              {{ recognizeResult.confidence_level || '-' }}
            </el-tag>
          </el-descriptions-item>
          <el-descriptions-item :label="$t('auto.List.19.6fcade')">
            <el-tag :type="recognizeResult.method === 'llm' ? 'danger' : 'info'" size="small" effect="plain">
              {{ recognizeResult.method === 'llm' ? 'LLM 模型' : '规则匹配' }}
            </el-tag>
          </el-descriptions-item>
          <el-descriptions-item v-if="recognizeResult.platform" :label="$t('auto.List.6.498128')">
            {{ recognizeResult.platform }}
          </el-descriptions-item>
          <el-descriptions-item v-if="recognizeResult.latency_ms" :label="$t('auto.List.20.39f137')">
            {{ recognizeResult.latency_ms }} ms
          </el-descriptions-item>
        </el-descriptions>

        <template v-if="recognizeResult.entities && Object.keys(recognizeResult.entities).length">
          <el-divider content-position="left">{{ $t('auto.List.49.d3d9ad') }}</el-divider>
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
          <el-table-column :label="$t('auto.List.21.932ccc')" min-width="200" show-overflow-tooltip>
            <template #default="{ row }">{{ row.message || row.raw_text }}</template>
          </el-table-column>
          <el-table-column :label="$t('auto.List.22.764aca')" width="120">
            <template #default="{ row }">
              <el-tag :type="getIntentTagType(row.intent_type)" size="small">
                {{ row.intent_name || intentNameMap[row.intent_type] || row.intent_type }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column :label="$t('auto.List.17.8f3599')" width="100">
            <template #default="{ row }">{{ Math.round((row.confidence || 0) * 100) }}%</template>
          </el-table-column>
          <el-table-column :label="$t('auto.List.23.ea340b')" width="80">
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
          <div class="stat-label">{{ $t('auto.List.50.2616b1') }}</div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card shadow="hover" class="stat-card stat-purchase">
          <div class="stat-value">{{ statsData.purchase || 0 }}</div>
          <div class="stat-label">{{ $t('auto.List.51.ce092e') }}</div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card shadow="hover" class="stat-card stat-inquiry">
          <div class="stat-value">{{ statsData.price_inquiry || 0 }}</div>
          <div class="stat-label">{{ $t('auto.List.52.be446c') }}</div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card shadow="hover" class="stat-card stat-objection">
          <div class="stat-value">{{ objectionTotal }}</div>
          <div class="stat-label">{{ $t('auto.List.53.ec0ea0') }}</div>
        </el-card>
      </el-col>
    </el-row>

    <el-row :gutter="16">
      
      <el-col :span="10">
        <el-card shadow="never">
          <template #header>
            <div class="card-header">
              <span>{{ $t('auto.List.54.549a89') }}</span>
              <el-button link type="primary" size="small" @click="loadStats">{{ $t('auto.List.55.694fc5') }}</el-button>
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
            <el-empty v-else :description="$t('auto.List.24.a3b437')" :image-size="80" />
          </div>
        </el-card>
      </el-col>

      
      <el-col :span="14">
        <el-card shadow="never">
          <template #header>
            <div class="card-header">
              <span>{{ $t('auto.List.56.2ca96e') }}</span>
              <div class="header-controls">
                <el-select v-model="filterIntentType" :placeholder="$t('auto.List.25.c212e9')" clearable size="small" style="width: 160px" @change="onFilterChange">
                  <el-option v-for="d in intentDict" :key="d.type" :label="d.name" :value="d.type" />
                </el-select>
                <el-button size="small" type="primary" @click="loadRecent">{{ $t('auto.List.57.bee912') }}</el-button>
              </div>
            </div>
          </template>
          <el-table :data="recentList" v-loading="recentLoading" stripe size="small" style="width: 100%">
            <el-table-column :label="$t('auto.List.26.b87b77')" min-width="180" show-overflow-tooltip>
              <template #default="{ row }">{{ row.message || row.raw_text }}</template>
            </el-table-column>
            <el-table-column :label="$t('auto.List.27.c6d069')" width="120">
              <template #default="{ row }">
                <el-tag :type="getIntentTagType(row.intent_type)" size="small">
                  {{ intentNameMap[row.intent_type] || row.intent_type }}
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column :label="$t('auto.List.17.8f3599')" width="130">
              <template #default="{ row }">
                <el-progress
                  :percentage="Math.round((row.confidence || 0) * 100)"
                  :color="getConfidenceColor(row.confidence_level)"
                  :stroke-width="12"
                  :text-inside="true"
                />
              </template>
            </el-table-column>
            <el-table-column :label="$t('auto.List.6.498128')" width="90">
              <template #default="{ row }">
                <el-tag v-if="row.platform" size="small" effect="plain">{{ getChannelLabel(row.platform) }}</el-tag>
                <span v-else>-</span>
              </template>
            </el-table-column>
            <el-table-column :label="$t('auto.List.28.19fcb9')" width="150">
              <template #default="{ row }">{{ formatTime(row.created_at) }}</template>
            </el-table-column>
            <template #empty>
              <el-empty :description="$t('auto.List.29.b9db84')" :image-size="80" />
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
          <span><el-icon><MagicStick /></el-icon>{{ $t('auto.List.58.902eda') }}</span>
          <el-button link type="primary" @click="fillFineExample">{{ $t('填入示例') }}</el-button>
        </div>
      </template>
      <el-form :model="fineForm" label-width="90px" inline>
        <el-form-item :label="$t('auto.List.14.80ded0')">
          <el-input v-model="fineForm.message" type="textarea" :rows="2" :placeholder="$t('auto.List.30.c8be06')" style="width: 460px" />
        </el-form-item>
        <el-form-item :label="$t('auto.List.4.3f61de')">
          <el-input v-model="fineForm.customer_id" :placeholder="$t('auto.List.31.c20cba')" clearable style="width: 180px" />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" :loading="fineLoading" @click="runRecognizeFine">
            <el-icon><MagicStick /></el-icon> 精细识别
          </el-button>
        </el-form-item>
      </el-form>
      <template v-if="fineResult">
        <el-divider content-position="left">{{ $t('auto.List.59.20d4ea') }}</el-divider>
        <el-descriptions :column="2" border>
          <el-descriptions-item :label="$t('auto.List.32.d196c7')">
            <el-tag type="primary" effect="dark">{{ fineResult.major || '-' }}</el-tag>
          </el-descriptions-item>
          <el-descriptions-item :label="$t('auto.List.33.815e21')">
            <el-tag effect="plain">{{ fineResult.minor || '-' }}</el-tag>
          </el-descriptions-item>
          <el-descriptions-item :label="$t('auto.List.17.8f3599')">
            <el-progress
              :percentage="Math.round((fineResult.confidence || 0) * 100)"
              :color="(fineResult.confidence || 0) >= 0.7 ? '#10B981' : '#F59E0B'"
              :stroke-width="14"
              style="width: 220px"
            />
          </el-descriptions-item>
          <el-descriptions-item :label="$t('auto.List.23.ea340b')">
            <el-tag :type="getFineMethodType(fineResult.method)" size="small">
              {{ getFineMethodLabel(fineResult.method) }}
            </el-tag>
          </el-descriptions-item>
          <el-descriptions-item v-if="fineResult.llm_model" :label="$t('auto.List.34.c82b0d')">
            <el-tag size="small" effect="plain">{{ fineResult.llm_model }}</el-tag>
          </el-descriptions-item>
          <el-descriptions-item v-if="fineResult.latency_ms" :label="$t('auto.List.20.39f137')">
            {{ fineResult.latency_ms }} ms
          </el-descriptions-item>
          <el-descriptions-item v-if="fineResult.reasoning" :label="$t('auto.List.35.633fb3')" :span="2">
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
            <el-select v-model="filterFineMajor" :placeholder="$t('auto.List.36.30ad2a')" clearable size="small" style="width: 160px" @change="loadFineLogs">
              <el-option v-for="m in fineMajorOptions" :key="m" :label="m" :value="m" />
            </el-select>
            <el-button size="small" type="primary" @click="loadFineLogs">{{ $t('auto.List.57.bee912') }}</el-button>
          </div>
        </div>
      </template>
      <el-table :data="fineLogs" v-loading="fineLogsLoading" stripe size="small">
        <el-table-column :label="$t('auto.List.26.b87b77')" min-width="200" show-overflow-tooltip>
          <template #default="{ row }">{{ row.message || '-' }}</template>
        </el-table-column>
        <el-table-column :label="$t('auto.List.32.d196c7')" width="120">
          <template #default="{ row }">
            <el-tag type="primary" size="small" effect="plain">{{ row.intent_major || '-' }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column :label="$t('auto.List.33.815e21')" width="160" show-overflow-tooltip>
          <template #default="{ row }">{{ row.intent_minor || '-' }}</template>
        </el-table-column>
        <el-table-column :label="$t('auto.List.17.8f3599')" width="100">
          <template #default="{ row }">{{ Math.round((row.confidence || 0) * 100) }}%</template>
        </el-table-column>
        <el-table-column :label="$t('auto.List.23.ea340b')" width="80">
          <template #default="{ row }">
            <el-tag :type="getFineMethodType(row.method)" size="small" effect="plain">
              {{ getFineMethodLabel(row.method) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column :label="$t('auto.List.20.39f137')" width="80">
          <template #default="{ row }">{{ row.latency_ms || 0 }}ms</template>
        </el-table-column>
        <el-table-column :label="$t('auto.List.28.19fcb9')" width="150">
          <template #default="{ row }">{{ formatTime(row.timestamp || row.created_at) }}</template>
        </el-table-column>
        <template #empty>
          <el-empty :description="$t('auto.List.37.04c99b')" :image-size="80" />
        </template>
      </el-table>
    </el-card>

    
    <el-card shadow="never" class="dict-card">
      <template #header>
        <div class="card-header">
          <span>意图词典（{{ intentDict.length }} 类）</span>
          <el-input v-model="dictKeyword" :placeholder="$t('auto.List.38.f12a40')" clearable size="small" style="width: 220px" />
        </div>
      </template>
      <el-table :data="filteredDict" stripe size="small">
        <el-table-column :label="$t('auto.List.22.764aca')" width="130">
          <template #default="{ row }">
            <el-tag :type="getIntentTagType(row.type)" size="small">{{ row.name }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="type" :label="$t('auto.List.39.45235f')" width="180" />
        <el-table-column prop="description" :label="$t('auto.List.40.f411d0')" min-width="180" show-overflow-tooltip />
        <el-table-column :label="$t('auto.List.41.9699a5')" min-width="260">
          <template #default="{ row }">
            <el-tag v-for="kw in (row.keywords || []).slice(0, 6)" :key="kw" size="small" effect="plain" style="margin: 2px">
              {{ kw }}
            </el-tag>
            <span v-if="(row.keywords || []).length > 6" class="more-tag">等 {{ row.keywords.length }} 个</span>
          </template>
        </el-table-column>
        <el-table-column :label="$t('auto.List.42.1a63ac')" min-width="220" show-overflow-tooltip>
          <template #default="{ row }">
            {{ (row.examples || []).join(' / ') }}
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <!-- 词表配置对话框 -->
    <el-dialog v-model="keywordOverrideVisible" :title="$t('auto.List.43.24690f')" width="720px">
      <el-alert
        :title="$t('auto.List.44.38a2e6')"
        type="info"
        :closable="false"
        show-icon
        style="margin-bottom: 16px"
      />
      <div v-if="keywordOverrideLoading" style="text-align: center; padding: 40px">
        <el-icon class="is-loading" :size="32"><Loading /></el-icon>
        <p style="margin-top: 12px; color: #909399">{{ $t('auto.List.60.26b5bd') }}</p>
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
              :placeholder="$t('auto.List.45.5b5c82')"
              size="small"
              style="flex: 1"
              @keyup.enter="addKeyword(intent.type)"
            />
            <el-button size="small" @click="addKeyword(intent.type)">{{ $t('auto.List.61.b58c75') }}</el-button>
          </div>
        </div>
      </div>
      <template #footer>
        <el-button @click="keywordOverrideVisible = false">{{ $t('auto.List.62.625fb2') }}</el-button>
        <el-button type="primary" :loading="keywordOverrideSaving" @click="saveKeywordOverride">{{ $t('auto.List.63.be5fbb') }}</el-button>
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
const RECOGNIZE_METHOD = { llm: 'LLM', rule: t('auto.List.64.b0fae0'), keyword: t('auto.List.41.9699a5'), bert: 'BERT', hybrid: t('auto.List.65.41ce13') };
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
const FINE_METHOD_LABEL = { rule: t('auto.List.64.b0fae0'), llm: 'LLM', hybrid: t('auto.List.66.2e5edb') }
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
    ElMessage.warning(t('auto.List.67.2c5ca9') + (e.message || t('auto.List.68.974e74')))
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
    ElMessage.success(`意图识别已${val ? t('auto.List.2.7854b5') : t('auto.List.3.b15d91')}，业务链路立即生效`)
  } catch (e) {
    intentEnabled.value = !val;
    ElMessage.error(t('auto.List.69.d0759f') + (e.message || t('auto.List.68.974e74')))
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
    ElMessage.error(t('auto.List.70.4c832b') + (e.message || t('auto.List.68.974e74')))
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
    ElMessage.success(t('auto.List.71.4e023c'))
    keywordOverrideVisible.value = false
  } catch (e) {
    ElMessage.error(t('auto.List.72.036fe2') + (e.message || t('auto.List.68.974e74')))
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
    ElMessage.error(t('auto.List.73.4f0b6e') + (e.message || t('auto.List.68.974e74')))
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
    ElMessage.error(t('auto.List.74.b32725') + (e.message || t('auto.List.68.974e74')))
  } finally {
    recentLoading.value = false
  }
}

const loadDict = async () => {
  try {
    const res = await intentApi.getDict()
    intentDict.value = Array.isArray(res) ? res : (res?.list || [])
  } catch (e) {
    ElMessage.error(t('auto.List.75.a48487') + (e.message || t('auto.List.68.974e74')))
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
    t('auto.List.76.ce57ca'),
    t('auto.List.77.ea54a8'),
    t('auto.List.78.164802'),
    t('auto.List.79.ed2ed6'),
    t('auto.List.80.c1a0cf')
  ]
  testForm.value.message = examples[Math.floor(Math.random() * examples.length)]
}

const runRecognize = async () => {
  if (!testForm.value.message || !testForm.value.message.trim()) {
    ElMessage.warning(i18n.global.t(t('auto.List.81.f6acdc')))
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
    ElMessage.success(i18n.global.t(t('auto.List.82.483f84')))
    loadRecent()
  } catch (e) {
    ElMessage.error(t('auto.List.83.a3ebed') + (e.message || t('auto.List.68.974e74')))
  } finally {
    recognizing.value = false
  }
}

const runBatchRecognize = async () => {
  const lines = (testForm.value.message || '').split('\n').map(s => s.trim()).filter(Boolean);
  if (lines.length === 0) {
    ElMessage.warning(i18n.global.t(t('auto.List.84.41c20d')))
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
    ElMessage.error(t('auto.List.85.09d30f') + (e.message || t('auto.List.68.974e74')))
  } finally {
    batchLoading.value = false
  }
}

const fillFineExample = () => {
  const examples = [
    t('auto.List.86.d9158d'),
    t('auto.List.87.770412'),
    t('auto.List.88.191b56'),
    t('auto.List.89.75f259'),
    t('auto.List.90.5ed1ee'),
    t('auto.List.91.e4dd37')
  ]
  fineForm.value.message = examples[Math.floor(Math.random() * examples.length)]
};

const runRecognizeFine = async () => {
  if (!fineForm.value.message || !fineForm.value.message.trim()) {
    ElMessage.warning(i18n.global.t(t('auto.List.81.f6acdc')))
    return
  }
  fineLoading.value = true
  fineResult.value = null
  try {
    const data = { message: fineForm.value.message }
    if (fineForm.value.customer_id) data.customer_id = fineForm.value.customer_id
    const res = await intentApi.recognizeFine(data)
    fineResult.value = res
    ElMessage.success(i18n.global.t(t('auto.List.92.489a81')))
    loadFineLogs()
  } catch (e) {
    ElMessage.error(t('auto.List.93.3569c4') + (e.message || t('auto.List.68.974e74')))
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
    ElMessage.error(t('auto.List.94.aa7ea4') + (e.message || t('auto.List.68.974e74')))
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
