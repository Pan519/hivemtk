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
            <el-form-item :label="$t('auto.Edit.a708d3ec')" prop="agent_code">
              <el-input v-model="form.agent_code" :placeholder="$t('auto.Edit.9b4fcd69')" :disabled="isEdit" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item :label="$t('auto.Edit.63ecdffc')" prop="name">
              <el-input v-model="form.name" :placeholder="$t('auto.Edit.2272d5c4')" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item :label="$t('auto.Edit.40d4da9a')" prop="agent_type">
              <el-select v-model="form.agent_type" style="width: 100%">
                <el-option :label="$t('auto.Edit.facb6bff')" value="sales" />
                <el-option :label="$t('auto.Edit.afaab1c7')" value="customer_service" />
                <el-option :label="$t('auto.Edit.bac24769')" value="hybrid" />
              </el-select>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item :label="$t('auto.Edit.28a44ceb')" prop="agent_mode">
              <el-select v-model="form.agent_mode" style="width: 100%">
                <el-option :label="$t('auto.Edit.635fa101')" value="passive" />
                <el-option :label="$t('auto.Edit.691800b3')" value="active" />
              </el-select>
              <div class="form-tip">{{ $t('auto.Edit.c5821b6a') }}</div>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item :label="$t('auto.Edit.673c2033')">
              <el-select
                v-model="form.asset_bundle_id"
                :placeholder="$t('auto.Edit.3fc8e54d')"
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
              <div class="form-tip">{{ $t('auto.Edit.34fb9abc') }}</div>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item :label="$t('auto.Edit.21aa156d')">
              <el-input v-model="form.avatar" :placeholder="$t('auto.Edit.1ff94e3a')" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item :label="$t('auto.Edit.573278b4')" prop="internal_language">
              <el-select v-model="form.internal_language" :placeholder="$t('auto.Edit.5c0be50b')" style="width: 100%">
                <el-option
                  v-for="lang in internalLanguageOptions"
                  :key="lang.value"
                  :label="lang.label"
                  :value="lang.value"
                />
              </el-select>
              <div class="form-tip">{{ $t('auto.Edit.ff8ccf25') }}</div>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item :label="$t('auto.Edit.129d0594')" prop="target_language">
              <el-select v-model="form.target_language" :placeholder="$t('auto.Edit.c91ca080')" clearable style="width: 100%">
                <el-option
                  v-for="lang in targetLanguageOptions"
                  :key="lang.value"
                  :label="lang.label"
                  :value="lang.value"
                />
              </el-select>
              <div class="form-tip">{{ $t('auto.Edit.7caf6cbc') }}</div>
            </el-form-item>
          </el-col>
          <el-col :span="24">
            <el-form-item :label="$t('auto.Edit.3bdd08ad')">
              <el-input v-model="form.description" type="textarea" :rows="2" :placeholder="$t('auto.Edit.8a648af0')" />
            </el-form-item>
          </el-col>
        </el-row>
      </el-card>

      
      <el-card shadow="never" class="section-card">
        <template #header>
          <div class="card-title"><el-icon><UserFilled /></el-icon>{{ $t('auto.Edit.d4ffddff') }}</div>
        </template>
        <el-form-item :label="$t('auto.Edit.5749e0cf')" prop="persona">
          <el-input
            v-model="form.persona"
            type="textarea"
            :rows="4"
            :placeholder="$t('auto.Edit.7d36b648')"
          />
        </el-form-item>
        <el-form-item :label="$t('auto.Edit.cbdb6bc9')">
          <el-input
            v-model="form.system_prompt"
            type="textarea"
            :rows="4"
            :placeholder="$t('auto.Edit.2d9aef22')"
          />
        </el-form-item>
        <el-form-item :label="$t('auto.Edit.45eddfac')">
          <el-input v-model="form.greeting" :placeholder="$t('auto.Edit.55c6ec1e')" />
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
        <el-form-item :label="$t('auto.Edit.7f2c76c6')">
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
            :placeholder="$t('auto.Edit.76bfbd15')"
            style="width: 100%"
            :render-after-expand="false"
            ::empty-text="$t('auto.Edit.9e348ebf')"
          />
          <div class="form-tip">
            树形结构按 RAG / FAQ / SOP 三个顶层分组。挂载知识库后，对话中自动应用该 KB 下所有文档 / FAQ / SOP 模板。
            旧的 faq_entry_ids / sop_template_ids / rag_product_ids 字段已废弃但仍保留兼容。
          </div>
        </el-form-item>
        <el-row :gutter="20">
          <el-col :span="8">
            <el-form-item :label="$t('auto.Edit.57277094')">
              <el-select
                v-model="form.rag_product_ids"
                multiple
                filterable
                clearable
                :placeholder="$t('auto.Edit.270bdca7')"
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
            <el-form-item :label="$t('auto.Edit.e93654a7')">
              <el-select
                v-model="form.faq_entry_ids"
                multiple
                filterable
                clearable
                :placeholder="$t('auto.Edit.270bdca7')"
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
            <el-form-item :label="$t('auto.Edit.79e06e2b')">
              <el-select
                v-model="form.sop_template_ids"
                multiple
                filterable
                clearable
                :placeholder="$t('auto.Edit.270bdca7')"
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
          <div class="card-title"><el-icon><Connection /></el-icon>{{ $t('auto.Edit.2932d0e1') }}</div>
        </template>
        <el-form-item :label="$t('auto.Edit.9a78fb76')">
          <el-select
            v-model="form.sop_ids"
            multiple
            filterable
            clearable
            :placeholder="$t('auto.Edit.6354b294')"
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
          <div class="card-title"><el-icon><ChatLineSquare /></el-icon>{{ $t('auto.Edit.6399cc72') }}</div>
        </template>
        <el-form-item :label="$t('auto.Edit.d49563bf')">
          <el-select
            v-model="form.script_library_ids"
            multiple
            filterable
            clearable
            :placeholder="$t('auto.Edit.57869ce4')"
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
          <div class="card-title"><el-icon><Cpu /></el-icon>{{ $t('auto.Edit.742a0835') }}</div>
        </template>
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item :label="$t('auto.Edit.248d0fb8')">
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
            <el-form-item :label="$t('auto.Edit.b20b409f')">
              <el-input-number v-model="form.max_tokens" :min="100" :max="8192" :step="100" style="width: 100%" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item :label="$t('auto.Edit.c9bf0b88')">
              <el-slider v-model="form.temperature" :min="0" :max="2" :step="0.1" show-input />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="Top P">
              <el-slider v-model="form.top_p" :min="0" :max="1" :step="0.05" show-input />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item :label="$t('auto.Edit.b8f0c714')">
              <el-slider v-model="form.frequency_penalty" :min="-2" :max="2" :step="0.1" show-input />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item :label="$t('auto.Edit.ca5ae49a')">
              <el-slider v-model="form.presence_penalty" :min="-2" :max="2" :step="0.1" show-input />
            </el-form-item>
          </el-col>
        </el-row>
      </el-card>

      
      <el-card shadow="never" class="section-card">
        <template #header>
          <div class="card-title"><el-icon><Open /></el-icon>{{ $t('auto.Edit.8a79d656') }}</div>
        </template>
        <el-row :gutter="20">
          <el-col :span="8">
            <el-form-item :label="$t('auto.Edit.9a39a697')">
              <el-switch v-model="form.enable_rag" />
            </el-form-item>
          </el-col>
          <el-col :span="8">
            <el-form-item :label="$t('auto.Edit.b4545948')">
              <el-switch v-model="form.enable_script_match" />
            </el-form-item>
          </el-col>
          <el-col :span="8">
            <el-form-item :label="$t('auto.Edit.cf983876')">
              <el-switch v-model="form.enable_humanize_polish" />
            </el-form-item>
          </el-col>
          <el-col :span="8">
            <el-form-item :label="$t('auto.Edit.8d8c65e0')">
              <el-switch v-model="form.enable_content_audit" />
            </el-form-item>
          </el-col>
          <el-col :span="8">
            <el-form-item :label="$t('auto.Edit.655e227e')">
              <el-switch v-model="form.enable_playbook" />
            </el-form-item>
          </el-col>
        </el-row>
      </el-card>

      
      <el-card shadow="never" class="section-card">
        <template #header>
          <div class="card-title"><el-icon><Setting /></el-icon>{{ $t('auto.Edit.4bcc788b') }}</div>
        </template>
        <el-row :gutter="20">
          <el-col :span="8">
            <el-form-item label="RAG TopK">
              <el-input-number v-model="form.rag_top_k" :min="1" :max="20" :step="1" style="width: 100%" />
            </el-form-item>
          </el-col>
          <el-col :span="8">
            <el-form-item :label="$t('auto.Edit.bd493e1b')">
              <el-slider v-model="form.confidence_threshold" :min="0" :max="1" :step="0.05" show-input />
            </el-form-item>
          </el-col>
          <el-col :span="8">
            <el-form-item :label="$t('auto.Edit.755a7cad')">
              <el-input-number v-model="form.max_ai_consecutive" :min="1" :max="50" :step="1" style="width: 100%" />
            </el-form-item>
          </el-col>
        </el-row>
      </el-card>

      
      <div class="footer-actions">
        <el-button @click="goBack">{{ $t('auto.Edit.625fb26b') }}</el-button>
        <el-button type="primary" :loading="saving" @click="onSave">
          <el-icon><Check /></el-icon>
          {{ isEdit ? '保存更新' : '创建智能体' }}
        </el-button>
      </div>
    </el-form>

    
    <el-dialog v-model="testDialogVisible" :title="$t('auto.Edit.353e8df4')" width="760px" top="6vh">
      <el-form :model="testForm" label-width="90px">
        <el-form-item :label="$t('auto.Edit.72fbb9b5')">
          <el-input :model-value="form.name || form.agent_code" disabled />
        </el-form-item>
        <el-form-item :label="$t('auto.Edit.1c2ecf70')">
          <el-input v-model="testForm.customer_id" :placeholder="$t('auto.Edit.7bfb0d81')" />
        </el-form-item>
        <el-form-item :label="$t('auto.Edit.b87b7756')" required>
          <el-input
            v-model="testForm.message"
            type="textarea"
            :rows="3"
            :placeholder="$t('auto.Edit.cc633c28')"
          />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" :loading="testing" @click="runTest">
            <el-icon><VideoPlay /></el-icon>
            执行测试
          </el-button>
          <el-button @click="testResult = null">{{ $t('auto.Edit.7e61acdd') }}</el-button>
        </el-form-item>
      </el-form>

      <template v-if="testResult">
        <el-divider content-position="left">{{ $t('auto.Edit.c6c9da2b') }}</el-divider>
        <el-descriptions :column="2" border size="small">
          <el-descriptions-item :label="$t('auto.Edit.fab08d35')" :span="2">
            <div class="reply-box">{{ testResult.reply || '-' }}</div>
          </el-descriptions-item>
          <el-descriptions-item :label="$t('auto.Edit.248d0fb8')">{{ testResult.llm_model || '-' }}</el-descriptions-item>
          <el-descriptions-item :label="$t('auto.Edit.7f62c087')">{{ testResult.latency_ms || 0 }}</el-descriptions-item>
          <el-descriptions-item :label="$t('auto.Edit.8f5f4e35')">{{ testResult.cost_tokens || 0 }}</el-descriptions-item>
          <el-descriptions-item :label="$t('auto.Edit.39e5fbd5')">
            <el-tag :type="testResult.transferred_to_human ? 'warning' : 'success'" size="small">
              {{ testResult.transferred_to_human ? '是' : '否' }}
            </el-tag>
            <span v-if="testResult.transfer_reason" class="transfer-reason">（{{ testResult.transfer_reason }}）</span>
          </el-descriptions-item>
        </el-descriptions>

        <div class="chain-title">{{ $t('auto.Edit.cda843b5') }}</div>
        <el-table :data="testResult.steps || []" stripe size="small" border>
          <el-table-column type="index" label="#" width="50" align="center" />
          <el-table-column prop="step" :label="$t('auto.Edit.52b36576')" min-width="160" show-overflow-tooltip />
          <el-table-column :label="$t('auto.Edit.3fea7ca7')" width="90" align="center">
            <template #default="{ row }">
              <el-tag :type="getStepStatusType(row.status)" size="small">{{ row.status || '-' }}</el-tag>
            </template>
          </el-table-column>
          <el-table-column :label="$t('auto.Edit.7f62c087')" width="100" align="center">
            <template #default="{ row }">{{ row.latency_ms || 0 }}</template>
          </el-table-column>
          <el-table-column prop="detail" :label="$t('auto.Edit.f26225bd')" min-width="200" show-overflow-tooltip />
          <el-table-column prop="error" :label="$t('auto.Edit.7030ff64')" min-width="160" show-overflow-tooltip>
            <template #default="{ row }">
              <span v-if="row.error" class="error-text">{{ row.error }}</span>
              <span v-else>-</span>
            </template>
          </el-table-column>
        </el-table>
      </template>
      <el-empty v-else-if="!testing" :description="$t('auto.Edit.5910a5ed')" :image-size="80" />
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
    { required: true, message: i18n.global.t(t('auto.Edit.f817e581')), trigger: 'blur' },
    { pattern: /^[a-zA-Z0-9_-]+$/, message: i18n.global.t(t('auto.Edit.2777ef66')), trigger: 'blur' }
  ],
  name: [
    { required: true, message: i18n.global.t(t('auto.Edit.2272d5c4')), trigger: 'blur' }
  ],
  agent_type: [
    { required: true, message: i18n.global.t(t('auto.Edit.0251affb')), trigger: 'change' }
  ],
  persona: [
    { required: true, message: i18n.global.t(t('auto.Edit.984951e8')), trigger: 'blur' }
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
    console.warn(t('auto.Edit.f3373253'), e?.message);
  }
};

const loadSopOptions = async () => {
  try {
    const res = await sopApi.list({ page: 1, page_size: 100 })
    sopOptions.value = res?.list || res?.items || []
  } catch (e) {
    console.warn(t('auto.Edit.c6f24ae2'), e?.message)
  }
}

const loadScriptOptions = async () => {
  try {
    const res = await getScriptTemplateList({ page: 1, page_size: 100 })
    scriptOptions.value = res?.list || res?.items || []
  } catch (e) {
    console.warn(t('auto.Edit.2b96b20e'), e?.message)
  }
}

const loadAssetBundles = async () => {
  try {
    const res = await listBundles({ page: 1, page_size: 200, status: 'active' })
    assetBundleOptions.value = res?.list || res?.items || []
  } catch (e) {
    console.warn(t('auto.Edit.b1e8a62a'), e?.message)
  }
}

const loadFaqOptions = async () => {
  try {
    const res = await faqApi.list({ page: 1, page_size: 200 })
    faqOptions.value = res?.list || []
  } catch (e) {
    console.warn(t('auto.Edit.db4e4c72'), e?.message)
  }
}

const loadSopTemplateOptions = async () => {
  try {
    const res = await sopTemplateApi.list({ page: 1, page_size: 100 })
    sopTemplateOptions.value = res?.list || []
  } catch (e) {
    console.warn(t('auto.Edit.3a1ee2d2'), e?.message)
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
      tree.push({ id: 'kb-empty', label: t('auto.Edit.6f714237'), isLeaf: true, disabled: true });
    }
    kbTree.value = tree
    return tree
  } catch (e) {
    console.warn(t('auto.Edit.fc9815d2'), e?.message)
    kbTree.value = [{ id: 'kb-error', label: t('auto.Edit.c4fa9b74') + (e?.message || t('auto.Edit.974e7484')), isLeaf: true, disabled: true }]
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
    console.warn(t('auto.Edit.0d5a9dd4'), e?.message)
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
    ElMessage.error(t('auto.Edit.6ad039e9') + (e.message || t('auto.Edit.974e7484')))
    goBack()
  } finally {
    pageLoading.value = false
  }
};

const onSave = async () => {
  if (!formRef.value) return
  await formRef.value.validate(async (valid) => {
    if (!valid) {
      ElMessage.warning(i18n.global.t(t('auto.Edit.341de0e4')))
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
          console.warn(t('auto.Edit.3a622cb6'), kbErr?.message)
        }
        ElMessage.success(i18n.global.t(t('auto.Edit.8baf0b46')))
      } else {
        const created = await createAgent(data)
        const newId = created?.id || created?.ID;
        if (newId && form.kb_ids?.length) {
          try {
            await replaceAgentKBs(newId, form.kb_ids.map(String))
          } catch (kbErr) {
            console.warn(t('auto.Edit.047a0ffb'), kbErr?.message)
          }
        }
        ElMessage.success(i18n.global.t(t('auto.Edit.0cca28cd')))
      }
      goBack()
    } catch (e) {
      ElMessage.error(t('auto.Edit.78d235cc') + (e.message || t('auto.Edit.974e7484')))
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
    ElMessage.warning(i18n.global.t(t('auto.Edit.573b391d')))
    return
  }
  if (!isEdit.value) {
    ElMessage.warning(i18n.global.t(t('auto.Edit.f4521c3e')))
    return
  }
  testForm.value = { customer_id: '', message: '' }
  testResult.value = null
  testDialogVisible.value = true
};

const runTest = async () => {
  if (!testForm.value.message || !testForm.value.message.trim()) {
    ElMessage.warning(i18n.global.t(t('auto.Edit.fa967bd0')))
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
    ElMessage.success(i18n.global.t(t('auto.Edit.2afc9835')))
  } catch (e) {
    ElMessage.error(t('auto.Edit.4b3cb9fe') + (e.message || t('auto.Edit.974e7484')))
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
