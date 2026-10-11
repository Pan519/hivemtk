<template>
  <div class="lead-mining-page">
    <el-tabs v-model="activeTab">
      <el-tab-pane :label="$t('auto.Index.e67e0d7b')" name="setting">
        <el-card shadow="never">
          <el-form :model="form" label-width="120px" class="setting-form">
            <el-form-item :label="$t('auto.Index.9d4058e1')">
              <el-switch v-model="form.enabled" />
              <span class="hint">{{ $t('auto.Index.6e34ec5a') }}</span>
            </el-form-item>

            <el-form-item :label="$t('auto.Index.7d5c5d72')">
              <el-checkbox-group v-model="form.channels">
                <el-checkbox v-for="c in channelOptions" :key="c.value" :label="c.value">{{ c.label }}</el-checkbox>
              </el-checkbox-group>
              <div class="hint">{{ $t('auto.Index.40f49f67') }}</div>
            </el-form-item>

            <el-form-item :label="$t('auto.Index.9699a50e')">
              <el-input
                v-model="keywordsText"
                type="textarea"
                :rows="3"
                :placeholder="$t('auto.Index.1b0dfa74')"
              />
              <div class="hint">{{ $t('auto.Index.adb19e76') }}</div>
            </el-form-item>

            <el-form-item :label="$t('auto.Index.2f1ae528')">
              <el-input
                v-model="tagsText"
                type="textarea"
                :rows="2"
                :placeholder="$t('auto.Index.75a18217')"
              />
            </el-form-item>

            <el-form-item :label="$t('auto.Index.12e04d29')">
              <el-input
                v-model="form.requirement"
                type="textarea"
                :rows="4"
                :placeholder="$t('auto.Index.46a3ee8f')"
              />
            </el-form-item>

            <el-form-item :label="$t('auto.Index.bb3526e8')">
              <el-input-number v-model="form.min_intent_score" :min="0" :max="100" :step="5" />
              <span class="hint">{{ $t('auto.Index.84f741a8') }}</span>
            </el-form-item>

            <el-form-item :label="$t('auto.Index.92702b9a')">
              <el-input v-model="form.model" :placeholder="$t('auto.Index.832d7c0f')" style="max-width: 320px" />
            </el-form-item>

            <el-form-item>
              <el-button type="primary" :loading="saving" @click="onSave">{{ $t('auto.Index.4323be72') }}</el-button>
              <el-button @click="onReset">{{ $t('auto.Index.4b9c3271') }}</el-button>
            </el-form-item>
          </el-form>
        </el-card>
      </el-tab-pane>

      <el-tab-pane :label="$t('auto.Index.1e76b415')" name="leads">
        <el-card shadow="never">
          <div class="toolbar">
            <el-input
              v-model="kw"
              :placeholder="$t('auto.Index.feb0a9f8')"
              clearable
              style="width: 240px"
              @input="loadLeads"
            />
            <el-button type="primary" @click="loadLeads">{{ $t('auto.Index.694fc5ef') }}</el-button>
          </div>
          <el-table :data="leads" v-loading="loading" stripe style="margin-top: 12px">
            <el-table-column prop="name" :label="$t('auto.Index.d7ec2d3f')" min-width="160" />
            <el-table-column prop="desc" :label="$t('auto.Index.0de6682b')" min-width="320" show-overflow-tooltip />
            <el-table-column :label="$t('auto.Index.c2900498')" width="100">
              <template #default="{ row }">
                <el-tag :type="row.intent_score >= 70 ? 'danger' : (row.intent_score >= 40 ? 'warning' : 'info')">
                  {{ row.intent_score }}
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column :label="$t('auto.Index.55290e55')" width="90">
              <template #default="{ row }">
                <el-tag v-if="row.is_opportunity === 1 || row.is_opportunity === true" type="success">{{ $t('auto.Index.2601e297') }}</el-tag>
                <span v-else>-</span>
              </template>
            </el-table-column>
            <el-table-column prop="one_id" :label="$t('auto.Index.d409158f')" min-width="160" show-overflow-tooltip />
            <el-table-column :label="$t('auto.Index.2fe5f3d8')" width="170">
              <template #default="{ row }">{{ formatTime(row.create_time) }}</template>
            </el-table-column>
          </el-table>
          <el-pagination
            v-if="total > pageSize"
            layout="prev, pager, next"
            :total="total"
            :page-size="pageSize"
            @current-change="onPage"
            style="margin-top: 12px"
          />
        </el-card>
      </el-tab-pane>
    </el-tabs>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { ElMessage } from 'element-plus'
import { getLeadMiningConfig, saveLeadMiningConfig } from '@/api/leadMining'
import { clueApi } from '@/api/clue'

const { t } = useI18n()

const activeTab = ref('setting')
const saving = ref(false)
const keywordsText = ref('')
const tagsText = ref('')
const form = reactive({
  enabled: false,
  channels: [],
  requirement: '',
  min_intent_score: 50,
  model: ''
})

const channelOptions = [
  { value: 'telegram', label: 'Telegram' },
  { value: 'whatsapp', label: 'WhatsApp' },
  { value: 'douyin', label: t('auto.Index.277a984f') },
  { value: 'xiaohongshu', label: t('auto.Index.ae26bb04') },
  { value: 'wecom', label: t('auto.Index.ff17b9f9') },
  { value: 'wechat', label: t('auto.Index.cfbf6f4c') },
  { value: 'xianyu', label: t('auto.Index.9d4865ba') },
  { value: 'kuaishou', label: t('auto.Index.f32c0a2a') },
  { value: 'tiktok', label: 'TikTok' },
  { value: 'instagram', label: 'Instagram' }
]

function parseList(text) {
  return (text || '')
    .split(/[\n,，、]/)
    .map((s) => s.trim())
    .filter(Boolean)
}

async function loadConfig() {
  try {
    const res = await getLeadMiningConfig()
    // 拦截器已拆包：config 直接回配置对象，按信封取 .data 只会读到 undefined
    const data = res?.data?.data || res?.data || res || {}
    form.enabled = !!data.enabled
    form.channels = data.channels || []
    form.requirement = data.requirement || ''
    form.min_intent_score = data.min_intent_score || 50
    form.model = data.model || ''
    keywordsText.value = (data.keywords || []).join(', ')
    tagsText.value = (data.tags || []).join(', ')
  } catch (e) {
    ElMessage.error(t('auto.Index.2f367553') + (e.message || e))
  }
}

async function onSave() {
  saving.value = true
  const payload = {
    enabled: form.enabled,
    channels: form.channels,
    requirement: form.requirement,
    min_intent_score: form.min_intent_score,
    model: form.model,
    keywords: parseList(keywordsText.value),
    tags: parseList(tagsText.value)
  }
  try {
    await saveLeadMiningConfig(payload)
    ElMessage.success(t('auto.Index.3b108349'))
  } catch (e) {
    ElMessage.error(t('auto.Index.78d235cc') + (e.message || e))
  } finally {
    saving.value = false
  }
}

function onReset() {
  loadConfig()
}

const leads = ref([]);
const loading = ref(false)
const total = ref(0)
const pageSize = ref(20)
const page = ref(1)
const kw = ref('')

function formatTime(t) {
  if (!t) return '-'
  const d = new Date(t * 1000)
  if (isNaN(d.getTime())) return String(t)
  const p = (n) => String(n).padStart(2, '0')
  return `${d.getFullYear()}-${p(d.getMonth() + 1)}-${p(d.getDate())} ${p(d.getHours())}:${p(d.getMinutes())}`
}

async function loadLeads() {
  loading.value = true
  try {
    const res = await clueApi.list(page.value, pageSize.value)
    const envelope = res?.data ?? res
    const data = envelope?.data ?? envelope
    const list = data?.list ?? (Array.isArray(data) ? data : [])
    let filtered = list.filter((c) => Number(c.type) === 8);
    if (kw.value) {
      const k = kw.value.toLowerCase()
      filtered = filtered.filter(
        (c) => (c.name || '').toLowerCase().includes(k) || (c.desc || '').toLowerCase().includes(k)
      )
    }
    leads.value = filtered
    total.value = data?.total ?? list.length
  } catch (e) {
    ElMessage.error(t('auto.Index.a2a651dd') + (e.message || e))
  } finally {
    loading.value = false
  }
}

function onPage(p) {
  page.value = p
  loadLeads()
}

onMounted(() => {
  loadConfig()
  loadLeads()
})
</script>

<style scoped>
.lead-mining-page {
  padding: 16px;
}
.setting-form {
  max-width: 760px;
}
.hint {
  color: #909399;
  font-size: 12px;
  margin-left: 10px;
}
.toolbar {
  display: flex;
  gap: 12px;
  align-items: center;
}
</style>
