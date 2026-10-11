<template>
  <div class="geo-page">

    <div class="p-4">
    <el-row :gutter="16" class="mb-4">
      <el-col :span="8">
        <el-card>
          <el-statistic :title="$t('auto.PlatformPublish.a3b56572')" :value="availablePlatforms.length" />
        </el-card>
      </el-col>
      <el-col :span="8">
        <el-card>
          <el-statistic :title="$t('auto.PlatformPublish.23811cbc')" :value="articles.length" />
        </el-card>
      </el-col>
      <el-col :span="8">
        <el-card>
          <el-statistic :title="$t('auto.PlatformPublish.abddeeb7')" :value="recordsToday" />
        </el-card>
      </el-col>
    </el-row>

    
    <el-card class="mb-4">
      <template #header>
        <div class="flex items-center justify-between">
          <span class="font-bold">{{ $t('auto.PlatformPublish.e63f99fc') }}</span>
          <el-button size="small" type="primary" @click="loadAll">{{ $t('auto.PlatformPublish.694fc5ef') }}</el-button>
        </div>
      </template>
      <el-table :data="availablePlatforms" v-loading="platformsLoading" size="small">
        <el-table-column prop="display_name" :label="$t('auto.PlatformPublish.4981287f')" width="160">
          <template #default="{ row }">
            <a v-if="row.url" :href="row.url" target="_blank" class="platform-link">{{ row.display_name }}</a>
            <span v-else>{{ row.display_name }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="name" :label="$t('auto.PlatformPublish.45235f6c')" width="160" />
        <el-table-column :label="$t('auto.PlatformPublish.16bd32d2')" width="110">
          <template #default="{ row }">
            <el-tag v-if="row.capability === 'real_api'" type="success" size="small">{{ $t('auto.PlatformPublish.144d3805') }}</el-tag>
            <el-tag v-else-if="row.capability === 'cookie_gray'" type="warning" size="small">Cookie</el-tag>
            <el-tag v-else type="info" size="small">{{ row.capability || '未知' }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column :label="$t('auto.PlatformPublish.b33c7279')" width="110">
          <template #default="{ row }">
            <span class="auth-type">{{ authLabel(row.auth_type) }}</span>
          </template>
        </el-table-column>
        <el-table-column :label="$t('auto.PlatformPublish.3fea7ca7')" width="130">
          <template #default="{ row }">
            <el-tag v-if="isConfigured(row.name)" type="success" size="small">{{ $t('auto.PlatformPublish.da208e9c') }}</el-tag>
            <el-tag v-else type="info" size="small">{{ $t('auto.PlatformPublish.f6c72d76') }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column :label="$t('auto.PlatformPublish.f411d0f1')" min-width="360">
          <template #default="{ row }">
            <span v-if="row.name && row.name.includes('github')" class="hint">
              访问 <a href="https://github.com/settings/tokens" target="_blank">https://github.com/settings/tokens</a>{{ $t('auto.PlatformPublish.57a0d75d') }}<code>GITHUB_TOKEN</code>
            </span>
            <span v-else-if="row.name === 'medium'" class="hint">
              访问 <a href="https://medium.com/me/settings/security" target="_blank">https://medium.com/me/settings/security</a>{{ $t('auto.PlatformPublish.e203cd8e') }}<code>MEDIUM_TOKEN</code>
            </span>
            <span v-else-if="row.name === 'wordpress'" class="hint">
              WordPress 后台 → 用户 → 应用密码 → 生成专用密码，填入环境变量 <code>WORDPRESS_APPPASS</code>（用户名+密码）
            </span>
            <span v-else-if="row.capability === 'real_api'" class="hint">{{ $t('auto.PlatformPublish.65c68c84') }}</span>
            <span v-else class="hint">—</span>
          </template>
        </el-table-column>
        <el-table-column :label="$t('auto.PlatformPublish.2b6bc0f2')" width="200">
          <template #default="{ row }">
            <el-button size="small" @click="openAccountDialog(row)">{{ $t('auto.PlatformPublish.4e401048') }}</el-button>
            <el-button size="small" type="primary" plain :disabled="!isConfigured(row.name)" @click="onTestPublish(row)">{{ $t('auto.PlatformPublish.b5c54691') }}</el-button>
          </template>
        </el-table-column>
      </el-table>
      <div v-if="filteredCount === 0" class="empty-hint">
        暂无可用平台。需要在后端配置 API Token 后可发布。
      </div>
    </el-card>

    
    <el-card class="mb-4">
      <template #header><span class="font-bold">{{ $t('auto.PlatformPublish.51b2c299') }}</span></template>
      <div class="flex items-center gap-4 mb-4 flex-wrap">
        <el-select v-model="selectedArticle" :placeholder="$t('auto.PlatformPublish.b72fad55')" size="small" style="width:260px">
          <el-option v-for="a in articles" :key="a.id" :label="a.title || a.id" :value="a.id" />
        </el-select>
        <el-select v-model="selectedPlatforms" multiple :placeholder="$t('auto.PlatformPublish.e2c01ffc')" size="small" style="min-width:260px">
          <el-option v-for="p in realApiPlatforms" :key="p.name" ::label="$t('auto.PlatformPublish.25e49195')" :value="p.name" />
        </el-select>
        <el-button size="small" type="primary" :disabled="!selectedArticle || !selectedPlatforms.length" :loading="pipelineRunning" @click="runPipeline">{{ $t('auto.PlatformPublish.25de77cb') }}</el-button>
      </div>
      <el-table v-if="pipelineStatuses.length" :data="pipelineStatuses" size="small">
        <el-table-column prop="platform" :label="$t('auto.PlatformPublish.4981287f')" width="160" />
        <el-table-column prop="status" :label="$t('auto.PlatformPublish.3fea7ca7')" width="120">
          <template #default="{ row }">
            <el-tag :type="statusType(row.status)" size="small">{{ row.status }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="message" :label="$t('auto.PlatformPublish.ff692f04')" min-width="200" show-overflow-tooltip />
        <el-table-column prop="url" :label="$t('auto.PlatformPublish.1a3b0804')" min-width="200">
          <template #default="{ row }">
            <a v-if="row.url" :href="row.url" target="_blank">{{ row.url }}</a>
            <span v-else>-</span>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    
    <el-card>
      <template #header><span class="font-bold">{{ $t('auto.PlatformPublish.c1534887') }}</span></template>
      <el-table :data="records" v-loading="recordsLoading" size="small">
        <el-table-column prop="article_title" :label="$t('auto.PlatformPublish.c75625dc')" min-width="160" show-overflow-tooltip />
        <el-table-column prop="platform" :label="$t('auto.PlatformPublish.4981287f')" width="160" />
        <el-table-column prop="status" :label="$t('auto.PlatformPublish.3fea7ca7')" width="100">
          <template #default="{ row }">
            <el-tag :type="statusType(row.status)" size="small">{{ row.status }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="error" :label="$t('auto.PlatformPublish.4604d502')" min-width="200" show-overflow-tooltip />
        <el-table-column prop="created_at" :label="$t('auto.PlatformPublish.19fcb9eb')" width="180" />
      </el-table>
      <div v-if="records.length === 0" class="empty-hint">{{ $t('auto.PlatformPublish.9b97e3c1') }}</div>
    </el-card>

    
    <el-dialog v-model="accountDialogVisible" :title="`配置 ${currentPlatform?.display_name} Token`" width="480px">
      <el-form :model="accountForm" label-width="120px">
        <el-form-item label="AppKey / UserID">
          <el-input v-model="accountForm.app_key" :placeholder="$t('auto.PlatformPublish.05a7fe59')" />
        </el-form-item>
        <el-form-item label="Secret / Token">
          <el-input v-model="accountForm.secret" type="password" show-password :placeholder="$t('auto.PlatformPublish.02b73d5b')" />
        </el-form-item>
        <el-form-item :label="$t('auto.PlatformPublish.05b5a8e5')">
          <el-input v-model="accountForm.extra" type="textarea" :rows="3" :placeholder="$t('auto.PlatformPublish.fcb91a13')" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="accountDialogVisible = false">{{ $t('auto.PlatformPublish.625fb26b') }}</el-button>
        <el-button type="primary" @click="saveAccount">{{ $t('auto.PlatformPublish.be5fbbe3') }}</el-button>
      </template>
    </el-dialog>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, reactive, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { ElMessage } from 'element-plus'
import { geoApi } from '@/api/geo.js'
import {
  listPlatforms as fetchPlatforms,
  listPlatformAccounts,
  savePlatformAccount,
  publishToPlatform,
  listPlatformPublishRecords
} from '@/api/geoPlatform.js'

const { t } = useI18n()

const allPlatforms = ref([]);
const platformsLoading = ref(false)
const accounts = ref({})
const articles = ref([])
const records = ref([])
const recordsLoading = ref(false)

const selectedArticle = ref(null)
const selectedPlatforms = ref([])
const pipelineRunning = ref(false)
const pipelineStatuses = ref([])

const accountDialogVisible = ref(false)
const currentPlatform = ref(null)
const accountForm = reactive({ app_key: '', secret: '', extra: '' })

const availablePlatforms = computed(() => {
  return allPlatforms.value.filter(p => {
    const cap = p.capability
    if (cap === 'stub') return false
    if (cap === 'cookie_gray' && !p.has_account) return false
    return true
  })
});

const hiddenCount = computed(() => allPlatforms.value.length - availablePlatforms.value.length)
const filteredCount = computed(() => availablePlatforms.value.length)

const realApiPlatforms = computed(() =>
  availablePlatforms.value.filter(p => p.capability === 'real_api')
)

const recordsToday = computed(() => {
  const t = new Date().toISOString().slice(0, 10)
  return records.value.filter(r => (r.created_at || '').startsWith(t)).length
})

const configuredPlatformNames = computed(() =>
  Object.keys(accounts.value).filter(k => accounts.value[k] && (accounts.value[k].app_key || accounts.value[k].configured))
)

const isConfigured = (name) => configuredPlatformNames.value.includes(name)

const authLabel = (t) => {
  const map = { token: 'API Token', oauth: 'OAuth', cookie: 'Cookie', xmlrpc: 'XML-RPC', custom: t('auto.PlatformPublish.f1d4ff50') }
  return map[t] || t || '-'
}

const statusType = (s) => {
  if (s === 'success' || s === 'done') return 'success'
  if (s === 'failed' || s === 'error') return 'danger'
  if (s === 'running' || s === 'pending') return 'warning'
  return 'info'
}

const loadAll = async () => {
  await Promise.all([loadPlatforms(), loadAccounts(), loadArticles(), loadRecords()])
}

const loadPlatforms = async () => {
  platformsLoading.value = true
  try {
    const data = await fetchPlatforms()
    allPlatforms.value = Array.isArray(data) ? data : []
  } catch {
    allPlatforms.value = []
  } finally {
    platformsLoading.value = false
  }
}

const loadAccounts = async () => {
  try {
    const data = await listPlatformAccounts()
    accounts.value = Array.isArray(data)
      ? data.reduce((acc, a) => { acc[a.platform || a.name || a.key] = a; return acc }, {})
      : (data || {})
  } catch { /* 忽略：清理/存储/恢复类 best-effort 操作 */ }
}

const loadArticles = async () => {
  try {
    const data = await geoApi.getArticleList({ page: 1, limit: 20 })
    articles.value = data?.list || data?.items || (Array.isArray(data) ? data : [])
  } catch { articles.value = [] }
}

const loadRecords = async () => {
  recordsLoading.value = true
  try {
    const data = await listPlatformPublishRecords(selectedArticle.value)
    records.value = Array.isArray(data) ? data : (data?.list || data?.items || [])
  } catch { records.value = [] }
  finally { recordsLoading.value = false }
}

const openAccountDialog = (p) => {
  currentPlatform.value = p
  const existing = accounts.value[p.name] || accounts.value[p.key] || {}
  accountForm.app_key = existing.app_key || existing.appKey || ''
  accountForm.secret = existing.secret || existing.token || ''
  accountForm.extra = typeof existing.extra === 'string' ? existing.extra : JSON.stringify(existing.extra || {})
  accountDialogVisible.value = true
}

const saveAccount = async () => {
  try {
    await savePlatformAccount({
      platform: currentPlatform.value.name || currentPlatform.value.key,
      app_key: accountForm.app_key,
      secret: accountForm.secret,
      extra: accountForm.extra
    })
    ElMessage.success(t('auto.PlatformPublish.3b108349'))
    accountDialogVisible.value = false
    await loadAccounts()
  } catch (e) {
    ElMessage.error(e?.message || t('auto.PlatformPublish.6de920b4'))
  }
}

const onTestPublish = async (p) => {
  if (!selectedArticle.value && !articles.value.length) {
    ElMessage.warning(t('auto.PlatformPublish.3c44892c'))
    return
  }
  const articleId = selectedArticle.value || articles.value[0]?.id
  const platformKey = p.name || p.key
  try {
    await publishToPlatform(articleId, platformKey, {})
    ElMessage.success(`${p.display_name} 测试发布成功`)
  } catch (e) {
    ElMessage.error(`${p.display_name} 测试失败：${e?.message || e}`)
  }
}

const runPipeline = async () => {
  if (!selectedArticle.value || !selectedPlatforms.value.length) {
    ElMessage.warning(t('auto.PlatformPublish.ed9c22a7'))
    return
  }
  pipelineRunning.value = true
  pipelineStatuses.value = selectedPlatforms.value.map(p => ({ platform: p, status: 'running', message: t('auto.PlatformPublish.27d8a2fa') }))
  try {
    for (const p of selectedPlatforms.value) {
      const idx = pipelineStatuses.value.findIndex(s => s.platform === p)
      try {
        await publishToPlatform(selectedArticle.value, p, {})
        if (idx >= 0) pipelineStatuses.value[idx] = { platform: p, status: 'success', message: t('auto.PlatformPublish.2fff079b') }
      } catch (err) {
        if (idx >= 0) pipelineStatuses.value[idx] = { platform: p, status: 'failed', message: err?.message || t('auto.PlatformPublish.250688d7') }
      }
    }
    try {
      const recordsData = await listPlatformPublishRecords(selectedArticle.value)
      const list = Array.isArray(recordsData) ? recordsData : (recordsData?.list || recordsData?.items || [])
      if (list.length) {
        pipelineStatuses.value = list.map(r => ({
          platform: r.platform,
          status: r.status || 'unknown',
          message: r.message || r.error || '',
          url: r.url
        }))
      }
    } catch { /* 忽略：清理/存储/恢复类 best-effort 操作 */ }
    pipelineRunning.value = false
    ElMessage.success('Pipeline 执行完成')
    loadRecords()
  } catch (e) {
    ElMessage.error('Pipeline 执行失败：' + (e?.message || e))
    pipelineRunning.value = false
  }
}

onMounted(loadAll)
</script>

<style lang="scss" scoped>
.platform-link {
  color: var(--el-color-primary);
  text-decoration: none;
  &:hover { text-decoration: underline; }
}
.auth-type {
  font-size: 12px;
  color: var(--el-text-color-secondary);
}
.hint {
  font-size: 12px;
  color: var(--el-text-color-secondary);
}
.empty-hint {
  text-align: center;
  padding: 24px;
  color: var(--el-text-color-secondary);
  font-size: 13px;
}
.hidden-hint {
  margin-top: 12px;
}
</style>
