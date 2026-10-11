<template>
  <div class="api-token-page">
    <el-alert
      ::title="$t('auto.ApiToken.97c6dcb7')"
      type="info"
      :closable="false"
      show-icon
    >
      <template #default>
        为外部系统(自有 CRM/ERP/Helpdesk)创建 API Token,可通过 <code>/api/knowledge-merchant/external/import</code> 推送文档。
        创建时明文仅显示一次,请立即保存;之后系统只保存哈希。
      </template>
    </el-alert>

    <el-row :gutter="16" style="margin-top: 16px">
      
      <el-col :span="16">
        <el-card>
          <template #header>
            <div class="card-header">
              <span>{{ $t('auto.ApiToken.94460b5a') }}</span>
              <div>
                <el-select v-model="filterProductId" :placeholder="$t('auto.ApiToken.6437ee21')" clearable style="width: 180px" @change="loadTokens">
                  <el-option v-for="p in productList" :key="p.id" :label="p.name" :value="p.id" />
                </el-select>
                <el-button :icon="Refresh" @click="loadTokens" style="margin-left: 8px">{{ $t('auto.ApiToken.694fc5ef') }}</el-button>
              </div>
            </div>
          </template>
          <el-table :data="tokens" v-loading="loading" stripe>
            <el-table-column prop="id" label="ID" width="80" />
            <el-table-column prop="name" :label="$t('auto.ApiToken.d7ec2d3f')" min-width="140" show-overflow-tooltip />
            <el-table-column :label="$t('auto.ApiToken.a015434e')" width="160" show-overflow-tooltip>
              <template #default="{ row }">
                {{ getProductName(row.product_id) }}
              </template>
            </el-table-column>
            <el-table-column :label="$t('auto.ApiToken.ba6e91fa')" width="160">
              <template #default="{ row }">
                <el-tag v-for="s in parseScopes(row.scopes)" :key="s" size="small" style="margin-right: 4px">
                  {{ scopeLabel(s) }}
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column :label="$t('auto.ApiToken.3fea7ca7')" width="100">
              <template #default="{ row }">
                <el-tag v-if="row.enabled === 1" type="success">{{ $t('auto.ApiToken.7854b52a') }}</el-tag>
                <el-tag v-else type="danger">{{ $t('auto.ApiToken.579886e6') }}</el-tag>
              </template>
            </el-table-column>
            <el-table-column :label="$t('auto.ApiToken.d389aac2')" width="140">
              <template #default="{ row }">
                <div>调用 {{ row.use_count || 0 }} 次</div>
                <el-text v-if="row.last_used_at" size="small" type="info">最后: {{ formatDate(row.last_used_at) }}</el-text>
              </template>
            </el-table-column>
            <el-table-column :label="$t('auto.ApiToken.1fa23f4d')" width="170">
              <template #default="{ row }">
                <span v-if="row.expires_at">{{ formatDate(row.expires_at) }}</span>
                <span v-else style="color: #909399">{{ $t('auto.ApiToken.718ab6a0') }}</span>
              </template>
            </el-table-column>
            <el-table-column prop="created_by" :label="$t('auto.ApiToken.95a43eaa')" width="100" />
            <el-table-column :label="$t('auto.ApiToken.2b6bc0f2')" width="140" fixed="right">
              <template #default="{ row }">
                <el-button v-if="row.enabled === 1" link type="danger" size="small" @click="handleRevoke(row)">{{ $t('auto.ApiToken.e87fd913') }}</el-button>
                <el-text v-else type="info" size="small">{{ $t('auto.ApiToken.579886e6') }}</el-text>
              </template>
            </el-table-column>
          </el-table>
        </el-card>
      </el-col>

      
      <el-col :span="8">
        <el-card>
          <template #header>
            <span>{{ $t('auto.ApiToken.710c40cb') }}</span>
          </template>
          <el-form :model="form" :rules="rules" ref="formRef" label-width="100px">
            <el-form-item :label="$t('auto.ApiToken.d7ec2d3f')" prop="name">
              <el-input v-model="form.name" :placeholder="$t('auto.ApiToken.818bf17e')" />
            </el-form-item>
            <el-form-item :label="$t('auto.ApiToken.a015434e')" prop="product_id">
              <el-select v-model="form.product_id" :placeholder="$t('auto.ApiToken.cbf1cfff')" style="width: 100%">
                <el-option v-for="p in productList" :key="p.id" :label="p.name" :value="p.id" />
              </el-select>
            </el-form-item>
            <el-form-item :label="$t('auto.ApiToken.ba6e91fa')">
              <el-checkbox-group v-model="form.scopes">
                <el-checkbox label="read">{{ $t('auto.ApiToken.85541bd9') }}</el-checkbox>
                <el-checkbox label="write">{{ $t('auto.ApiToken.a5230ad6') }}</el-checkbox>
              </el-checkbox-group>
              <div class="form-hint">{{ $t('auto.ApiToken.1f17e0c8') }}</div>
            </el-form-item>
            <el-form-item :label="$t('auto.ApiToken.1fa23f4d')">
              <el-date-picker
                v-model="form.expires_at"
                type="datetime"
                :placeholder="$t('auto.ApiToken.81d3b24d')"
                style="width: 100%"
                value-format="YYYY-MM-DDTHH:mm:ssZ"
                format="YYYY-MM-DD HH:mm"
              />
            </el-form-item>
            <el-form-item>
              <el-button type="primary" :loading="creating" @click="handleCreate" style="width: 100%">{{ $t('auto.ApiToken.d9ac9228') }}</el-button>
            </el-form-item>
          </el-form>
        </el-card>

        
        <el-card style="margin-top: 16px">
          <template #header>
            <span>{{ $t('auto.ApiToken.9df2c688') }}</span>
          </template>
          <pre class="code-block">
# 1. 推送文档(JSON)
curl -X POST \
  'http://your-domain/api/knowledge-merchant/external/import' \
  -H 'X-Knowledge-Token: kbg_xxxx' \
  -H 'Content-Type: application/json' \
  -d '{
    "source": "custom",
    "product_id": "your-product",
    "items": [
      {"title": "FAQ1", "content": "..."}
    ]
  }'

# 2. 同步返回
curl -X POST \
  'http://your-domain/api/knowledge-merchant/external/import?sync=true' \
  -H 'X-Knowledge-Token: kbg_xxxx' \
  -H 'Content-Type: application/json' \
  -d '{"source":"custom","product_id":"...","items":[]}'</pre>
        </el-card>
      </el-col>
    </el-row>

    
    <el-dialog v-model="showTokenDialog" :title="$t('auto.ApiToken.f6a7301d')" width="640px" :close-on-click-modal="false">
      <el-alert type="success" :closable="false" show-icon>
        <template #default>
          请立即复制并保存明文 Token,关闭后无法再次查看。
        </template>
      </el-alert>
      <el-form label-width="80px" style="margin-top: 16px">
        <el-form-item :label="$t('auto.ApiToken.d7ec2d3f')">
          <el-text>{{ createdToken.name }}</el-text>
        </el-form-item>
        <el-form-item :label="$t('auto.ApiToken.a015434e')">
          <el-text>{{ getProductName(createdToken.product_id) }}</el-text>
        </el-form-item>
        <el-form-item label="Token">
          <el-input v-model="createdToken.token_plain" readonly>
            <template #append>
              <el-button :icon="CopyDocument" @click="copyToken">{{ $t('auto.ApiToken.79d3abe9') }}</el-button>
            </template>
          </el-input>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button type="primary" @click="showTokenDialog = false">{{ $t('auto.ApiToken.ac4dfef5') }}</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import i18n from '@/i18n'

import { ref, reactive, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Refresh, CopyDocument } from '@element-plus/icons-vue'
import { knowledgeMerchantAPI } from '@/api/knowledgeMerchant'
import { ragProductConfigAPI } from '@/api/ragProductConfig'

const { t } = useI18n()

const loading = ref(false)
const creating = ref(false)
const tokens = ref([])
const productList = ref([])
const filterProductId = ref('')

const formRef = ref(null)
const form = reactive({
  name: '',
  product_id: '',
  scopes: ['read', 'write'],
  expires_at: null
})

const rules = {
  name: [{ required: true, message: i18n.global.t(t('auto.ApiToken.06e2f88f')), trigger: 'blur' }],
  product_id: [{ required: true, message: i18n.global.t(t('auto.ApiToken.3b8e12b1')), trigger: 'change' }]
}

const showTokenDialog = ref(false)
const createdToken = ref({ name: '', product_id: 0, token_plain: '' })

const loadProducts = async () => {
  try {
    const res = await ragProductConfigAPI.listProducts()
    if (Array.isArray(res)) productList.value = res
    else if (res?.items) productList.value = res.items
  } catch (e) {
    console.error(t('auto.ApiToken.fcbbfe35'), e)
  }
}

const loadTokens = async () => {
  loading.value = true
  try {
    const res = await knowledgeMerchantAPI.listTokens({
      product_id: filterProductId.value
    })
    tokens.value = res?.items || res || []
  } catch (e) {
    ElMessage.error(t('auto.ApiToken.db3a28fd') + (e.message || ''))
  } finally {
    loading.value = false
  }
}

const getProductName = (id) => {
  const p = productList.value.find(p => p.id === id)
  return p ? p.name : `#${id}`
}

const parseScopes = (scopes) => {
  try {
    return JSON.parse(scopes)
  } catch {
    return []
  }
}

const scopeLabel = (s) => {
  return { read: t('auto.ApiToken.75769d1a'), write: t('auto.ApiToken.4d7dc6c5'), '*': t('auto.ApiToken.a8b0c204') }[s] || s
}

const handleCreate = async () => {
  if (!formRef.value) return
  const valid = await formRef.value.validate().catch(() => false)
  if (!valid) return
  creating.value = true
  try {
    const res = await knowledgeMerchantAPI.createToken({
      name: form.name,
      product_id: form.product_id,
      scopes: form.scopes,
      expires_at: form.expires_at || null
    })
    if (res) {
      createdToken.value = res
      showTokenDialog.value = true
      ElMessage.success(i18n.global.t('Token 已创建'))
      formRef.value.resetFields()
      form.scopes = ['read', 'write']
      form.expires_at = null
      loadTokens()
    }
  } catch (e) {
    ElMessage.error(t('auto.ApiToken.a2bddca3') + (e.message || ''))
  } finally {
    creating.value = false
  }
}

const handleRevoke = async (row) => {
  try {
    await ElMessageBox.confirm(
      `确认吊销 Token「${row.name}」吗?吊销后外部系统将无法继续使用。`,
      t('auto.ApiToken.a6ff25af'),
      { type: 'warning' }
    )
    await knowledgeMerchantAPI.revokeToken(row.id)
    ElMessage.success(i18n.global.t(t('auto.ApiToken.579886e6')))
    loadTokens()
  } catch (e) {
    if (e !== 'cancel') ElMessage.error(t('auto.ApiToken.cb3afa0c') + (e.message || ''))
  }
}

const copyToken = async () => {
  try {
    await navigator.clipboard.writeText(createdToken.value.token_plain)
    ElMessage.success(i18n.global.t(t('auto.ApiToken.4fb42e6e')))
  } catch {
    ElMessage.warning(i18n.global.t(t('auto.ApiToken.c39f1216')))
  }
}

const formatDate = (d) => d ? new Date(d).toLocaleString('zh-CN') : '-'

onMounted(async () => {
  await loadProducts()
  await loadTokens()
})
</script>

<style scoped lang="scss">
.api-token-page {
  padding: 0;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.form-hint {
  font-size: 12px;
  color: #909399;
  margin-top: 4px;
}

.code-block {
  background: #1e1e1e;
  color: #d4d4d4;
  padding: 12px;
  border-radius: 4px;
  font-size: 12px;
  line-height: 1.5;
  overflow-x: auto;
  white-space: pre-wrap;
  word-break: break-all;
}
</style>
