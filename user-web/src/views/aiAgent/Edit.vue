<template>
  <div class="ai-agent-edit-page">
    
    <el-card class="header-card" shadow="never">
      <div class="header-content">
        <div>
          <h2>{{ isEdit ? '编辑智能体' : '创建智能体' }}</h2>
          <p class="subtitle" v-if="isEdit">智能体ID：{{ agentId }} · 修改后保存生效</p>
          <p class="subtitle" v-else>{{ $t('填写智能体配置，创建后即可绑定到渠道或客服座席') }}</p>
        </div>
        <div class="header-actions">
          <el-button @click="goBack">{{ $t('返回列表') }}</el-button>
          <el-button type="primary" :loading="testing" @click="openTestDialog">
            <el-icon><VideoPlay /></el-icon>
            {{ $t('测试') }}
          </el-button>
        </div>
      </div>
    </el-card>

    <el-form
      ref="formRef"
      :model="form"
      :rules="rules"
      label-width="120px"
      v-loading="pageLoading"
    >
      
      <el-card shadow="never" class="section-card">
        <template #header>
          <div class="card-title"><el-icon><InfoFilled /></el-icon> {{ $t('基本信息') }}</div>
        </template>
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item :label="$t('auto.Edit.1.a708d3')" prop="agent_code">
              <el-input v-model="form.agent_code" :placeholder="$t('auto.Edit.2.9b4fcd')" :disabled="isEdit" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item :label="$t('auto.Edit.3.63ecdf')" prop="name">
              <el-input v-model="form.name" :placeholder="$t('auto.Edit.4.2272d5')" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item :label="$t('auto.Edit.5.40d4da')" prop="agent_type">
              <el-select v-model="form.agent_type" style="width: 100%">
                <el-option :label="$t('auto.Edit.6.facb6b')" value="sales" />
                <el-option :label="$t('auto.Edit.7.afaab1')" value="customer_service" />
                <el-option :label="$t('auto.Edit.8.bac247')" value="hybrid" />
              </el-select>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item :label="$t('auto.Edit.9.28a44c')" prop="agent_mode">
              <el-select v-model="form.agent_mode" style="width: 100%">
                <el-option :label="$t('auto.Edit.10.635fa1')" value="passive" />
                <el-option :label="$t('auto.Edit.11.691800')" value="active" />
              </el-select>
              <div class="form-tip">{{ $t('auto.Edit.66.c5821b') }}</div>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item :label="$t('auto.Edit.12.673c20')">
              <el-select
                v-model="form.asset_bundle_id"
                :placeholder="$t('auto.Edit.13.3fc8e5')"
                clearable
                filterable
                style="width: 100%"
              >
                <el-option
                  v-for="b in assetBundleOptions"
                  :key="b.asset_id"
                  :label="b.title || b.asset_id"
                  :value="b.asset_id"
                >
                  <span>{{ b.title || b.asset_id }}</span>
                  <span style="float: right; color: #8492a6; font-size: 12px">{{ b.industry || '' }}</span>
                </el-option>
              </el-select>
              <div class="form-tip">{{ $t('auto.Edit.67.34fb9a') }}</div>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item :label="$t('auto.Edit.14.21aa15')">
              <el-input v-model="form.avatar" :placeholder="$t('auto.Edit.15.1ff94e')" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item :label="$t('auto.Edit.16.573278')" prop="internal_language">
              <el-select v-model="form.internal_language" :placeholder="$t('auto.Edit.17.5c0be5')" style="width: 100%">
                <el-option
                  v-for="lang in internalLanguageOptions"
                  :key="lang.value"
                  :label="lang.label"
                  :value="lang.value"
                />
              </el-select>
              <div class="form-tip">{{ $t('auto.Edit.68.ff8ccf') }}</div>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item :label="$t('auto.Edit.18.129d05')" prop="target_language">
              <el-select v-model="form.target_language" :placeholder="$t('auto.Edit.19.c91ca0')" clearable style="width: 100%">
                <el-option
                  v-for="lang in targetLanguageOptions"
                  :key="lang.value"
                  :label="lang.label"
                  :value="lang.value"
                />
              </el-select>
              <div class="form-tip">{{ $t('auto.Edit.69.7caf6c') }}</div>
            </el-form-item>
          </el-col>
          <el-col :span="24">
            <el-form-item :label="$t('auto.Edit.20.3bdd08')">
              <el-input v-model="form.description" type="textarea" :rows="2" :placeholder="$t('auto.Edit.21.8a648a')" />
            </el-form-item>
          </el-col>
        </el-row>
      </el-card>

      
      <el-card shadow="never" class="section-card">
        <template #header>
          <div class="card-title"><el-icon><UserFilled /></el-icon>{{ $t('auto.Edit.70.d4ffdd') }}</div>
        </template>
        <el-form-item :label="$t('auto.Edit.22.5749e0')" prop="persona">
          <el-input
            v-model="form.persona"
            type="textarea"
            :rows="4"
            :placeholder="$t('auto.Edit.23.7d36b6')"
          />
        </el-form-item>
        <el-form-item :label="$t('auto.Edit.24.cbdb6b')">
          <el-input
            v-model="form.system_prompt"
            type="textarea"
            :rows="4"
            :placeholder="$t('auto.Edit.25.2d9aef')"
          />
        </el-form-item>
        <el-form-item :label="$t('auto.Edit.26.45eddf')">
          <el-input v-model="form.greeting" :placeholder="$t('auto.Edit.27.55c6ec')" />
        </el-form-item>
      </el-card>

      
      <el-card shadow="never" class="section-card">
        <template #header>
          <div class="card-title">
            <el-icon><Collection /></el-icon> 知识库挂载（树形多选）
            <el-button link type="primary" size="small" style="margin-left: auto" @click="loadKnowledgeBases" :loading="loadingKB">
              <el-icon><Refresh /></el-icon> 加载知识库
            </el-button>
          </div>
        </template>
        <el-form-item :label="$t('auto.Edit.28.7f2c76')">
          <el-tree-select
            v-model="form.kb_ids"
            :data="kbTree"
            :props="kbTreeProps"
            node-key="id"
            multiple
            show-checkbox
            check-strictly
            check-on-click-node
            collapse-tags
            collapse-tags-tooltip
            clearable
            filterable
            :default-expand-all="false"
            :placeholder="$t('auto.Edit.29.76bfbd')"
            style="width: 100%"
            :render-after-expand="false"
            ::empty-text="$t('auto.Edit.30.9e348e')"
          />
          <div class="form-tip">
            树形结构按 RAG / FAQ / SOP 三个顶层分组。挂载知识库后，对话中自动应用该 KB 下所有文档 / FAQ / SOP 模板。
            旧的 faq_entry_ids / sop_template_ids / rag_product_ids 字段已废弃但仍保留兼容。
          </div>
        </el-form-item>
        <el-row :gutter="20">
          <el-col :span="8">
            <el-form-item :label="$t('auto.Edit.31.572770')">
              <el-select
                v-model="form.rag_product_ids"
                multiple
                filterable
                clearable
                :placeholder="$t('auto.Edit.32.270bdc')"
                style="width: 100%"
              >
                <el-option
                  v-for="item in ragProductOptions"
                  :key="item.id"
                  :label="item.name"
                  :value="String(item.id)"
                />
              </el-select>
            </el-form-item>
          </el-col>
          <el-col :span="8">
            <el-form-item :label="$t('auto.Edit.33.e93654')">
              <el-select
                v-model="form.faq_entry_ids"
                multiple
                filterable
                clearable
                :placeholder="$t('auto.Edit.32.270bdc')"
                style="width: 100%"
              >
                <el-option
                  v-for="item in faqOptions"
                  :key="item.id"
                  :label="`#${item.id} ${item.question}`"
                  :value="String(item.id)"
                />
              </el-select>
            </el-form-item>
          </el-col>
          <el-col :span="8">
            <el-form-item :label="$t('auto.Edit.34.79e06e')">
              <el-select
                v-model="form.sop_template_ids"
                multiple
                filterable
                clearable
                :placeholder="$t('auto.Edit.32.270bdc')"
                style="width: 100%"
              >
                <el-option
                  v-for="item in sopTemplateOptions"
                  :key="item.id"
                  :label="`#${item.id} ${item.name}`"
                  :value="String(item.id)"
                />
              </el-select>
            </el-form-item>
          </el-col>
        </el-row>
      </el-card>

      
      <el-card shadow="never" class="section-card">
        <template #header>
          <div class="card-title"><el-icon><Connection /></el-icon>{{ $t('auto.Edit.71.2932d0') }}</div>
        </template>
        <el-form-item :label="$t('auto.Edit.35.9a78fb')">
          <el-select
            v-model="form.sop_ids"
            multiple
            filterable
            clearable
            :placeholder="$t('auto.Edit.36.6354b2')"
            style="width: 100%"
          >
            <el-option
              v-for="item in sopOptions"
              :key="item.id"
              :label="item.name"
              :value="String(item.id)"
            />
          </el-select>
        </el-form-item>
      </el-card>

      
      <el-card shadow="never" class="section-card">
        <template #header>
          <div class="card-title"><el-icon><ChatLineSquare /></el-icon>{{ $t('auto.Edit.72.6399cc') }}</div>
        </template>
        <el-form-item :label="$t('auto.Edit.37.d49563')">
          <el-select
            v-model="form.script_library_ids"
            multiple
            filterable
            clearable
            :placeholder="$t('auto.Edit.38.57869c')"
            style="width: 100%"
          >
            <el-option
              v-for="item in scriptOptions"
              :key="item.id"
              :label="item.title"
              :value="String(item.id)"
            />
          </el-select>
        </el-form-item>
      </el-card>

      
      <el-card shadow="never" class="section-card">
        <template #header>
          <div class="card-title"><el-icon><Cpu /></el-icon>{{ $t('auto.Edit.73.742a08') }}</div>
        </template>
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item :label="$t('auto.Edit.39.248d0f')">
              <el-select v-model="form.llm_model" style="width: 100%">
                <el-option label="gpt-4o-mini" value="gpt-4o-mini" />
                <el-option label="gpt-4o" value="gpt-4o" />
                <el-option label="gpt-4.1-mini" value="gpt-4.1-mini" />
                <el-option label="gpt-4.1" value="gpt-4.1" />
                <el-option label="claude-3-5-sonnet" value="claude-3-5-sonnet" />
                <el-option label="claude-3-5-haiku" value="claude-3-5-haiku" />
                <el-option label="deepseek-chat" value="deepseek-chat" />
                <el-option label="qwen-plus" value="qwen-plus" />
              </el-select>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item :label="$t('auto.Edit.40.b20b40')">
              <el-input-number v-model="form.max_tokens" :min="100" :max="8192" :step="100" style="width: 100%" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item :label="$t('auto.Edit.41.c9bf0b')">
              <el-slider v-model="form.temperature" :min="0" :max="2" :step="0.1" show-input />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="Top P">
              <el-slider v-model="form.top_p" :min="0" :max="1" :step="0.05" show-input />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item :label="$t('auto.Edit.42.b8f0c7')">
              <el-slider v-model="form.frequency_penalty" :min="-2" :max="2" :step="0.1" show-input />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item :label="$t('auto.Edit.43.ca5ae4')">
              <el-slider v-model="form.presence_penalty" :min="-2" :max="2" :step="0.1" show-input />
            </el-form-item>
          </el-col>
        </el-row>
      </el-card>

      
      <el-card shadow="never" class="section-card">
        <template #header>
          <div class="card-title"><el-icon><Open /></el-icon>{{ $t('auto.Edit.74.8a79d6') }}</div>
        </template>
        <el-row :gutter="20">
          <el-col :span="8">
            <el-form-item :label="$t('auto.Edit.44.9a39a6')">
              <el-switch v-model="form.enable_rag" />
            </el-form-item>
          </el-col>
          <el-col :span="8">
            <el-form-item :label="$t('auto.Edit.45.b45459')">
              <el-switch v-model="form.enable_script_match" />
            </el-form-item>
          </el-col>
          <el-col :span="8">
            <el-form-item :label="$t('auto.Edit.46.cf9838')">
              <el-switch v-model="form.enable_humanize_polish" />
            </el-form-item>
          </el-col>
          <el-col :span="8">
            <el-form-item :label="$t('auto.Edit.47.8d8c65')">
              <el-switch v-model="form.enable_content_audit" />
            </el-form-item>
          </el-col>
          <el-col :span="8">
            <el-form-item :label="$t('auto.Edit.48.655e22')">
              <el-switch v-model="form.enable_playbook" />
            </el-form-item>
          </el-col>
        </el-row>
      </el-card>

      
      <el-card shadow="never" class="section-card">
        <template #header>
          <div class="card-title"><el-icon><Setting /></el-icon>{{ $t('auto.Edit.75.4bcc78') }}</div>
        </template>
        <el-row :gutter="20">
          <el-col :span="8">
            <el-form-item label="RAG TopK">
              <el-input-number v-model="form.rag_top_k" :min="1" :max="20" :step="1" style="width: 100%" />
            </el-form-item>
          </el-col>
          <el-col :span="8">
            <el-form-item :label="$t('auto.Edit.49.bd493e')">
              <el-slider v-model="form.confidence_threshold" :min="0" :max="1" :step="0.05" show-input />
            </el-form-item>
          </el-col>
          <el-col :span="8">
            <el-form-item :label="$t('auto.Edit.50.755a7c')">
              <el-input-number v-model="form.max_ai_consecutive" :min="1" :max="50" :step="1" style="width: 100%" />
            </el-form-item>
          </el-col>
        </el-row>
      </el-card>

      
      <div class="footer-actions">
        <el-button @click="goBack">{{ $t('auto.Edit.76.625fb2') }}</el-button>
        <el-button type="primary" :loading="saving" @click="onSave">
          <el-icon><Check /></el-icon>
          {{ isEdit ? '保存更新' : '创建智能体' }}
        </el-button>
      </div>
    </el-form>

    
    <el-dialog v-model="testDialogVisible" :title="$t('auto.Edit.51.353e8d')" width="760px" top="6vh">
      <el-form :model="testForm" label-width="90px">
        <el-form-item :label="$t('auto.Edit.52.72fbb9')">
          <el-input :model-value="form.name || form.agent_code" disabled />
        </el-form-item>
        <el-form-item :label="$t('auto.Edit.53.1c2ecf')">
          <el-input v-model="testForm.customer_id" :placeholder="$t('auto.Edit.54.7bfb0d')" />
        </el-form-item>
        <el-form-item :label="$t('auto.Edit.55.b87b77')" required>
          <el-input
            v-model="testForm.message"
            type="textarea"
            :rows="3"
            :placeholder="$t('auto.Edit.56.cc633c')"
          />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" :loading="testing" @click="runTest">
            <el-icon><VideoPlay /></el-icon>
            执行测试
          </el-button>
          <el-button @click="testResult = null">{{ $t('auto.Edit.77.7e61ac') }}</el-button>
        </el-form-item>
      </el-form>

      <template v-if="testResult">
        <el-divider content-position="left">{{ $t('auto.Edit.78.c6c9da') }}</el-divider>
        <el-descriptions :column="2" border size="small">
          <el-descriptions-item :label="$t('auto.Edit.57.fab08d')" :span="2">
            <div class="reply-box">{{ testResult.reply || '-' }}</div>
          </el-descriptions-item>
          <el-descriptions-item :label="$t('auto.Edit.39.248d0f')">{{ testResult.llm_model || '-' }}</el-descriptions-item>
          <el-descriptions-item :label="$t('auto.Edit.58.7f62c0')">{{ testResult.latency_ms || 0 }}</el-descriptions-item>
          <el-descriptions-item :label="$t('auto.Edit.59.8f5f4e')">{{ testResult.cost_tokens || 0 }}</el-descriptions-item>
          <el-descriptions-item :label="$t('auto.Edit.60.39e5fb')">
            <el-tag :type="testResult.transferred_to_human ? 'warning' : 'success'" size="small">
              {{ testResult.transferred_to_human ? '是' : '否' }}
            </el-tag>
            <span v-if="testResult.transfer_reason" class="transfer-reason">（{{ testResult.transfer_reason }}）</span>
          </el-descriptions-item>
        </el-descriptions>

        <div class="chain-title">{{ $t('auto.Edit.79.cda843') }}</div>
        <el-table :data="testResult.steps || []" stripe size="small" border>
          <el-table-column type="index" label="#" width="50" align="center" />
          <el-table-column prop="step" :label="$t('auto.Edit.61.52b365')" min-width="160" show-overflow-tooltip />
          <el-table-column :label="$t('auto.Edit.62.3fea7c')" width="90" align="center">
            <template #default="{ row }">
              <el-tag :type="getStepStatusType(row.status)" size="small">{{ row.status || '-' }}</el-tag>
            </template>
          </el-table-column>
          <el-table-column :label="$t('auto.Edit.58.7f62c0')" width="100" align="center">
            <template #default="{ row }">{{ row.latency_ms || 0 }}</template>
          </el-table-column>
          <el-table-column prop="detail" :label="$t('auto.Edit.63.f26225')" min-width="200" show-overflow-tooltip />
          <el-table-column prop="error" :label="$t('auto.Edit.64.7030ff')" min-width="160" show-overflow-tooltip>
            <template #default="{ row }">
              <span v-if="row.error" class="error-text">{{ row.error }}</span>
              <span v-else>-</span>
            </template>
          </el-table-column>
        </el-table>
      </template>
      <el-empty v-else-if="!testing" :description="$t('auto.Edit.65.5910a5')" :image-size="80" />
    </el-dialog>
  </div>
</template>

<script setup>
import i18n from '@/i18n'

import { ref, reactive, computed, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import {
  InfoFilled, UserFilled, Collection, Connection, ChatLineSquare,
  Cpu, Open, Setting, Check, VideoPlay
} from '@element-plus/icons-vue'
import { createAgent, updateAgent, getAgent, testAgent } from '@/api/aiAgent.js'
import { ragProductConfigAPI } from '@/api/ragProductConfig'
import { sopApi } from '@/api/sopAgent.js'
import { getScriptTemplateList } from '@/api/scriptTemplate.js'
import { faqApi } from '@/api/faq'
import { sopTemplateApi } from '@/api/sopTemplate'
import { listByType as listKBByType } from '@/api/knowledgeBase'
import { listByAgent as listAgentKBs, replaceAgentKBs } from '@/api/agentKBBinding'
import { listBundles } from '@/api/assetBundle'
import { INTERNAL_LANGUAGE_OPTIONS, TARGET_LANGUAGE_OPTIONS } from '@/constants/languages'

const internalLanguageOptions = INTERNAL_LANGUAGE_OPTIONS
const targetLanguageOptions = TARGET_LANGUAGE_OPTIONS

const route = useRoute()
const router = useRouter()

const { t } = useI18n()

const formRef = ref();
const pageLoading = ref(false)
const saving = ref(false)
const agentId = computed(() => route.params.id)
const isEdit = computed(() => !!agentId.value)

const ragProductOptions = ref([]);
const sopOptions = ref([])
const scriptOptions = ref([])
const faqOptions = ref([])
const sopTemplateOptions = ref([])
const assetBundleOptions = ref([])

const kbTree = ref([]);
const loadingKB = ref(false)
const kbTreeProps = { value: 'id', label: 'label', children: 'children', isLeaf: 'isLeaf' }

const getDefaultForm = () => ({
  agent_code: '',
  name: '',
  description: '',
  avatar: '',
  agent_type: 'sales',
  agent_mode: 'passive',
  persona: '',
  system_prompt: '',
  greeting: '',
  asset_bundle_id: '',
  internal_language: 'zh',
  target_language: '',
  kb_ids: [],
  rag_product_ids: [],
  faq_entry_ids: [],
  sop_template_ids: [],
  sop_ids: [],
  script_library_ids: [],
  llm_model: 'gpt-4o-mini',
  temperature: 0.7,
  max_tokens: 800,
  top_p: 0.9,
  frequency_penalty: 0.5,
  presence_penalty: 0.5,
  enable_rag: true,
  enable_script_match: true,
  enable_humanize_polish: true,
  enable_content_audit: true,
  enable_playbook: true,
  rag_top_k: 3,
  confidence_threshold: 0.7,
  max_ai_consecutive: 5,
  status: 1
});

const form = reactive(getDefaultForm())

const rules = {
  agent_code: [
    { required: true, message: i18n.global.t(t('auto.Edit.80.f817e5')), trigger: 'blur' },
    { pattern: /^[a-zA-Z0-9_-]+$/, message: i18n.global.t(t('auto.Edit.81.2777ef')), trigger: 'blur' }
  ],
  name: [
    { required: true, message: i18n.global.t(t('auto.Edit.4.2272d5')), trigger: 'blur' }
  ],
  agent_type: [
    { required: true, message: i18n.global.t(t('auto.Edit.82.0251af')), trigger: 'change' }
  ],
  persona: [
    { required: true, message: i18n.global.t(t('auto.Edit.83.984951')), trigger: 'blur' }
  ]
};

const testDialogVisible = ref(false);
const testing = ref(false)
const testForm = ref({ customer_id: '', message: '' })
const testResult = ref(null)

const getStepStatusType = (status) => {
  if (status === 'ok') return 'success'
  if (status === 'fail') return 'danger'
  if (status === 'skip') return 'info'
  return 'info'
}

const loadRagProducts = async () => {
  try {
    const res = await ragProductConfigAPI.getRagProducts({ page: 1, page_size: 100 })
    ragProductOptions.value = res?.items || res?.list || []
  } catch (e) {
    console.warn(t('auto.Edit.84.f33732'), e?.message);
  }
};

const loadSopOptions = async () => {
  try {
    const res = await sopApi.list({ page: 1, page_size: 100 })
    sopOptions.value = res?.list || res?.items || []
  } catch (e) {
    console.warn(t('auto.Edit.85.c6f24a'), e?.message)
  }
}

const loadScriptOptions = async () => {
  try {
    const res = await getScriptTemplateList({ page: 1, page_size: 100 })
    scriptOptions.value = res?.list || res?.items || []
  } catch (e) {
    console.warn(t('auto.Edit.86.2b96b2'), e?.message)
  }
}

const loadAssetBundles = async () => {
  try {
    const res = await listBundles({ page: 1, page_size: 200, status: 'active' })
    assetBundleOptions.value = res?.list || res?.items || []
  } catch (e) {
    console.warn(t('auto.Edit.87.b1e8a6'), e?.message)
  }
}

const loadFaqOptions = async () => {
  try {
    const res = await faqApi.list({ page: 1, page_size: 200 })
    faqOptions.value = res?.list || []
  } catch (e) {
    console.warn(t('auto.Edit.88.db4e4c'), e?.message)
  }
}

const loadSopTemplateOptions = async () => {
  try {
    const res = await sopTemplateApi.list({ page: 1, page_size: 100 })
    sopTemplateOptions.value = res?.list || []
  } catch (e) {
    console.warn(t('auto.Edit.89.3a1ee2'), e?.message)
  }
}

const buildKBTree = async () => {
  loadingKB.value = true
  try {
    const [ragRes, faqRes, sopRes] = await Promise.all([
      listKBByType('rag', { page: 1, page_size: 200 }).catch(() => null),
      listKBByType('faq', { page: 1, page_size: 200 }).catch(() => null),
      listKBByType('sop', { page: 1, page_size: 200 }).catch(() => null)
    ])
    const extractList = (r) => (r?.list || r?.items || (Array.isArray(r) ? r : []))
    const ragList = extractList(ragRes)
    const faqList = extractList(faqRes)
    const sopList = extractList(sopRes)
    const buildChildren = (list) => list.map((it) => ({
      id: it.id,
      label: `[${it.kb_code || '#' + it.id}] ${it.name}`,
      isLeaf: true,
      kb_type: it.kb_type,
      raw: it
    }))
    const ragChildren = buildChildren(ragList)
    const faqChildren = buildChildren(faqList)
    const sopChildren = buildChildren(sopList)
    const tree = []
    if (ragChildren.length) {
      tree.push({ id: 'kb-type-rag', label: `RAG 文档库 (${ragChildren.length})`, children: ragChildren, isLeaf: false, disabled: true, kb_type: 'rag' })
    }
    if (faqChildren.length) {
      tree.push({ id: 'kb-type-faq', label: `FAQ 知识库 (${faqChildren.length})`, children: faqChildren, isLeaf: false, disabled: true, kb_type: 'faq' })
    }
    if (sopChildren.length) {
      tree.push({ id: 'kb-type-sop', label: `SOP 模板库 (${sopChildren.length})`, children: sopChildren, isLeaf: false, disabled: true, kb_type: 'sop' })
    }
    if (tree.length === 0) {
      tree.push({ id: 'kb-empty', label: t('auto.Edit.90.6f7142'), isLeaf: true, disabled: true });
    }
    kbTree.value = tree
    return tree
  } catch (e) {
    console.warn(t('auto.Edit.91.fc9815'), e?.message)
    kbTree.value = [{ id: 'kb-error', label: t('auto.Edit.92.c4fa9b') + (e?.message || t('auto.Edit.93.974e74')), isLeaf: true, disabled: true }]
    return []
  } finally {
    loadingKB.value = false
  }
};

const loadKnowledgeBases = async () => {
  await buildKBTree()
  if (isEdit.value && agentId.value) {
    await loadAgentKBBindings()
  }
}

const loadAgentKBBindings = async () => {
  try {
    const res = await listAgentKBs(agentId.value)
    const list = res?.list || res?.items || (Array.isArray(res) ? res : [])
    form.kb_ids = list.map((it) => it.kb_id ?? it.id)
  } catch (e) {
    console.warn(t('auto.Edit.94.0d5a9d'), e?.message)
    form.kb_ids = []
  }
};

const loadAgentDetail = async () => {
  if (!isEdit.value) return
  pageLoading.value = true
  try {
    const res = await getAgent(agentId.value)
    if (res) {
      Object.assign(form, getDefaultForm(), res)
      form.rag_product_ids = (res.rag_product_ids || []).map(String);
      form.faq_entry_ids = (res.faq_entry_ids || []).map(String)
      form.sop_template_ids = (res.sop_template_ids || []).map(String)
      form.sop_ids = (res.sop_ids || []).map(String)
      form.script_library_ids = (res.script_library_ids || []).map(String)
      form.kb_ids = (res.kb_ids || []).map(String)
      form.internal_language = res.internal_language || 'zh';
      form.target_language = res.target_language || ''
    }
  } catch (e) {
    ElMessage.error(t('auto.Edit.95.6ad039') + (e.message || t('auto.Edit.93.974e74')))
    goBack()
  } finally {
    pageLoading.value = false
  }
};

const onSave = async () => {
  if (!formRef.value) return
  await formRef.value.validate(async (valid) => {
    if (!valid) {
      ElMessage.warning(i18n.global.t(t('auto.Edit.96.341de0')))
      return
    }
    saving.value = true
    try {
      const data = {
        agent_code: form.agent_code,
        name: form.name,
        description: form.description,
        avatar: form.avatar,
        agent_type: form.agent_type,
        persona: form.persona,
        system_prompt: form.system_prompt,
        greeting: form.greeting,
        internal_language: form.internal_language || 'zh',
        target_language: form.target_language || '',
        kb_ids: form.kb_ids || [],
        rag_product_ids: form.rag_product_ids || [],
        faq_entry_ids: form.faq_entry_ids || [],
        sop_template_ids: form.sop_template_ids || [],
        sop_ids: form.sop_ids || [],
        script_library_ids: form.script_library_ids || [],
        llm_model: form.llm_model,
        temperature: form.temperature,
        max_tokens: form.max_tokens,
        top_p: form.top_p,
        frequency_penalty: form.frequency_penalty,
        presence_penalty: form.presence_penalty,
        enable_rag: form.enable_rag,
        enable_script_match: form.enable_script_match,
        enable_humanize_polish: form.enable_humanize_polish,
        enable_content_audit: form.enable_content_audit,
        enable_playbook: form.enable_playbook,
        rag_top_k: form.rag_top_k,
        confidence_threshold: form.confidence_threshold,
        max_ai_consecutive: form.max_ai_consecutive,
        status: form.status
      }
      if (isEdit.value) {
        await updateAgent(agentId.value, data)
        try {
          await replaceAgentKBs(agentId.value, form.kb_ids.map(String))
        } catch (kbErr) {
          console.warn(t('auto.Edit.97.3a622c'), kbErr?.message)
        }
        ElMessage.success(i18n.global.t(t('auto.Edit.98.8baf0b')))
      } else {
        const created = await createAgent(data)
        const newId = created?.id || created?.ID;
        if (newId && form.kb_ids?.length) {
          try {
            await replaceAgentKBs(newId, form.kb_ids.map(String))
          } catch (kbErr) {
            console.warn(t('auto.Edit.99.047a0f'), kbErr?.message)
          }
        }
        ElMessage.success(i18n.global.t(t('auto.Edit.100.0cca28')))
      }
      goBack()
    } catch (e) {
      ElMessage.error(t('auto.Edit.101.78d235') + (e.message || t('auto.Edit.93.974e74')))
    } finally {
      saving.value = false
    }
  })
};

const goBack = () => {
  router.push('/aiAgent/list')
}

const openTestDialog = () => {
  if (isEdit.value && !form.id) {
    ElMessage.warning(i18n.global.t(t('auto.Edit.102.573b39')))
    return
  }
  if (!isEdit.value) {
    ElMessage.warning(i18n.global.t(t('auto.Edit.103.f4521c')))
    return
  }
  testForm.value = { customer_id: '', message: '' }
  testResult.value = null
  testDialogVisible.value = true
};

const runTest = async () => {
  if (!testForm.value.message || !testForm.value.message.trim()) {
    ElMessage.warning(i18n.global.t(t('auto.Edit.104.fa967b')))
    return
  }
  testing.value = true
  testResult.value = null
  try {
    const data = {
      customer_id: testForm.value.customer_id || '',
      message: testForm.value.message
    }
    const res = await testAgent(agentId.value, data)
    testResult.value = res
    ElMessage.success(i18n.global.t(t('auto.Edit.105.2afc98')))
  } catch (e) {
    ElMessage.error(t('auto.Edit.106.4b3cb9') + (e.message || t('auto.Edit.93.974e74')))
  } finally {
    testing.value = false
  }
}

onMounted(async () => {
  await Promise.all([loadRagProducts(), loadSopOptions(), loadScriptOptions(), loadAssetBundles()])
  if (isEdit.value) {
    await loadAgentDetail()
  }
  buildKBTree().then(() => {
    if (isEdit.value) loadAgentKBBindings()
  });
});
</script>

<style scoped lang="scss">
.ai-agent-edit-page { padding: 20px; }

.header-card {
  margin-bottom: 16px;
  :deep(.el-card__body) { padding: 16px 20px; }
  .header-content {
    display: flex;
    justify-content: space-between;
    align-items: center;
    h2 { margin: 0 0 6px 0; font-size: 20px; }
    .subtitle { color: #909399; margin: 0; font-size: 13px; }
    .header-actions { display: flex; gap: 8px; }
  }
}

.section-card {
  margin-bottom: 16px;
  .card-title {
    display: flex;
    align-items: center;
    gap: 6px;
    font-weight: 600;
    font-size: 15px;
    color: #303133;
  }
}

.form-tip {
  font-size: 12px;
  color: #909399;
  line-height: 1.5;
  margin-top: 4px;
}

.footer-actions {
  display: flex;
  justify-content: center;
  gap: 12px;
  padding: 16px 0;
}

.reply-box {
  white-space: pre-wrap;
  word-break: break-word;
  background: #f5f7fa;
  padding: 10px 12px;
  border-radius: 4px;
  line-height: 1.6;
  max-height: 200px;
  overflow-y: auto;
}

.transfer-reason { color: #F59E0B; margin-left: 6px; }

.error-text { color: #EF4444; }

.chain-title {
  margin: 16px 0 8px 0;
  font-weight: 600;
  color: #303133;
  font-size: 14px;
}
</style>
