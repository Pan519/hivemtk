<template>
  <div class="obs-config">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>{{ $t('云存储配置') }}</span>
          <el-button type="primary" @click="handleCreate">{{ $t('新增配置') }}</el-button>
        </div>
      </template>

      <el-table :data="configList" v-loading="loading" border>
        <el-table-column prop="name" :label="$t('auto.ObsConfig.4fcad1c9')" width="150" />
        <el-table-column prop="provider" :label="$t('auto.ObsConfig.79f09f89')" width="120">
          <template #default="{ row }">
            <el-tag :type="getProviderType(row.provider)">
              {{ getProviderLabel(row.provider) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column :label="$t('auto.ObsConfig.e0e3f040')" min-width="240">
          <template #default="{ row }">
            <span v-if="row.provider === 'local'">
              {{ row.endpoint || './uploads' }}
            </span>
            <span v-else>
              {{ row.bucket || '-' }}<span v-if="row.endpoint"> @ {{ row.endpoint }}</span>
            </span>
          </template>
        </el-table-column>
        <el-table-column prop="domain" :label="$t('auto.ObsConfig.2139fb41')" min-width="180">
          <template #default="{ row }">
            <span v-if="row.domain">{{ row.domain }}</span>
            <span v-else class="muted">-</span>
          </template>
        </el-table-column>
        <el-table-column prop="is_default" :label="$t('auto.ObsConfig.18c63459')" width="80">
          <template #default="{ row }">
            <el-tag v-if="row.is_default" type="success">{{ $t('auto.ObsConfig.0a60ac8f') }}</el-tag>
            <span v-else>{{ $t('auto.ObsConfig.c9744f45') }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="status" :label="$t('auto.ObsConfig.3fea7ca7')" width="100">
          <template #default="{ row }">
            <el-tag :type="getStatusTagType(row.status)">
              {{ getStatusLabel(row.status) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column :label="$t('auto.ObsConfig.2b6bc0f2')" width="220" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" link @click="handleTest(row)">{{ $t('auto.ObsConfig.db06c78d') }}</el-button>
            <el-button type="primary" link @click="handleEdit(row)">{{ $t('auto.ObsConfig.95b351c8') }}</el-button>
            <el-button
              type="success"
              link
              @click="handleSetDefault(row)"
              v-if="!row.is_default"
            >
              设为默认
            </el-button>
            <el-button type="danger" link @click="handleDelete(row)">{{ $t('auto.ObsConfig.2f4aaddd') }}</el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    
    <el-dialog
      :title="dialogTitle"
      v-model="dialogVisible"
      width="640px"
      @close="handleDialogClose"
    >
      <el-form
        ref="formRef"
        :model="formData"
        :rules="dynamicRules"
        label-width="130px"
      >
        <el-form-item :label="$t('auto.ObsConfig.4fcad1c9')" prop="name">
          <el-input v-model="formData.name" :placeholder="$t('auto.ObsConfig.b984df04')" />
        </el-form-item>

        <el-form-item :label="$t('auto.ObsConfig.79f09f89')" prop="provider">
          <el-select v-model="formData.provider" :placeholder="$t('auto.ObsConfig.c52325b1')" @change="handleProviderChange" style="width:100%">
            <el-option :label="$t('auto.ObsConfig.5e2bb608')" value="local" />
            <el-option :label="$t('auto.ObsConfig.a4bd7c03')" value="aliyun" />
            <el-option :label="$t('auto.ObsConfig.3deda530')" value="qiniu" />
            <el-option :label="$t('auto.ObsConfig.b002b172')" value="tencent" />
            <el-option label="AWS S3" value="aws" />
          </el-select>
        </el-form-item>

        
        <template v-if="formData.provider !== 'local'">
          <el-form-item label="Access Key" prop="access_key">
            <el-input v-model="formData.access_key" :placeholder="$t('auto.ObsConfig.a9c49061')" show-password />
          </el-form-item>

          <el-form-item label="Secret Key" prop="secret_key">
            <el-input v-model="formData.secret_key" type="password" :placeholder="$t('auto.ObsConfig.4de5441d')" show-password />
          </el-form-item>

          <el-form-item :label="$t('auto.ObsConfig.bd5fcfda')" prop="bucket">
            <el-input v-model="formData.bucket" :placeholder="$t('auto.ObsConfig.e95f6949')" />
          </el-form-item>

          <el-form-item :label="$t('auto.ObsConfig.d7e47365')" prop="region" v-if="formData.provider === 'aws'">
            <el-input v-model="formData.region" :placeholder="$t('auto.ObsConfig.70644a48')" />
          </el-form-item>
        </template>

        
        <el-form-item :label="endpointLabel" prop="endpoint">
          <el-input v-model="formData.endpoint" :placeholder="endpointPlaceholder" />
          <div class="form-tip" v-if="formData.provider === 'local'">
            本地存储目录，默认 ./uploads（相对程序运行目录）
          </div>
          <div class="form-tip" v-else-if="formData.provider === 'aliyun'">
            https://xxx.oss-cn-hangzhou.aliyuncs.com
          </div>
          <div class="form-tip" v-else-if="formData.provider === 'qiniu'">
            http://xxx.bkt.clouddn.com
          </div>
          <div class="form-tip" v-else-if="formData.provider === 'tencent'">
            https://xxx.cos.ap-guangzhou.myqcloud.com
          </div>
          <div class="form-tip" v-else-if="formData.provider === 'aws'">
            https://s3.ap-southeast-1.amazonaws.com
          </div>
        </el-form-item>

        
        <el-form-item :label="domainLabel">
          <el-input v-model="formData.domain" :placeholder="domainPlaceholder" />
          <div class="form-tip" v-if="formData.provider === 'local'">
            本地文件对外访问的 URL 前缀，默认 /files
          </div>
          <div class="form-tip" v-else>
            可选：自定义 CDN 加速域名（留空则使用 endpoint 访问）
          </div>
        </el-form-item>

        <el-form-item :label="$t('auto.ObsConfig.475d55d8')">
          <el-input v-model="formData.path_prefix" :placeholder="$t('auto.ObsConfig.0cf3c8e8')" />
        </el-form-item>

        <el-row :gutter="16">
          <el-col :span="12">
            <el-form-item :label="$t('auto.ObsConfig.5bfe03ae')">
              <el-input-number v-model="formData.max_size" :min="10240" :step="10485760" />
              <div class="form-tip">{{ $t('auto.ObsConfig.72f7e41e') }}</div>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item :label="$t('auto.ObsConfig.37d0cc6b')">
              <el-input-number v-model="formData.max_count" :min="1" :step="100" />
              <div class="form-tip">{{ $t('auto.ObsConfig.59c0ea4b') }}</div>
            </el-form-item>
          </el-col>
        </el-row>

        <el-form-item :label="$t('auto.ObsConfig.3fea7ca7')">
          <el-radio-group v-model="formData.status">
            <el-radio label="active">{{ $t('auto.ObsConfig.7854b52a') }}</el-radio>
            <el-radio label="inactive">{{ $t('auto.ObsConfig.710ad08b') }}</el-radio>
          </el-radio-group>
        </el-form-item>
      </el-form>

      <template #footer>
        <el-button @click="dialogVisible = false">{{ $t('auto.ObsConfig.625fb26b') }}</el-button>
        <el-button type="primary" @click="handleSubmit">{{ $t('auto.ObsConfig.38cf16f2') }}</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import i18n from '@/i18n'
import { ref, reactive, computed, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { ElMessage, ElMessageBox } from 'element-plus'
import { getObsConfigList, createObsConfig, updateObsConfig, deleteObsConfig, testObsConnection, setDefaultObsConfig } from '@/api/obs'

const OBS_STATUS = { active: t('auto.ObsConfig.fd6e80f1'), disabled: t('auto.ObsConfig.710ad08b'), inactive: t('auto.ObsConfig.710ad08b'), error: t('auto.ObsConfig.7030ff64') };
const OBS_STATUS_TAG = { active: 'success', disabled: 'danger', inactive: 'danger', error: 'warning' }
const getStatusLabel = (s) => OBS_STATUS[s] || s || '-'
const getStatusTagType = (s) => OBS_STATUS_TAG[s] || ''

const { t } = useI18n()

const loading = ref(false)
const configList = ref([])
const dialogVisible = ref(false)
const formRef = ref(null)
const currentId = ref(null)

const formData = reactive({
  name: '',
  provider: '',
  access_key: '',
  secret_key: '',
  bucket: '',
  region: '',
  endpoint: '',
  domain: '',
  path_prefix: '',
  max_size: 104857600,
  max_count: 1000,
  status: 'active'
});

const dynamicRules = computed(() => {
  const rules = {
    name: [{ required: true, message: i18n.global.t(t('auto.ObsConfig.b984df04')), trigger: 'blur' }],
    provider: [{ required: true, message: i18n.global.t(t('auto.ObsConfig.c52325b1')), trigger: 'change' }]
  }
  if (formData.provider !== 'local') {
    rules.access_key = [{ required: true, message: i18n.global.t(t('auto.ObsConfig.a9c49061')), trigger: 'blur' }]
    rules.secret_key = [{ required: true, message: i18n.global.t(t('auto.ObsConfig.4de5441d')), trigger: 'blur' }]
    rules.bucket = [{ required: true, message: i18n.global.t(t('auto.ObsConfig.e95f6949')), trigger: 'blur' }]
  }
  if (formData.provider === 'aws') {
    rules.region = [{ required: true, message: i18n.global.t(t('auto.ObsConfig.09765cf3')), trigger: 'blur' }]
  }
  return rules
});

const dialogTitle = computed(() => (currentId.value ? t('auto.ObsConfig.5117bc6c') : t('auto.ObsConfig.80e2ca37')))

const endpointLabel = computed(() =>
  formData.provider === 'local' ? t('auto.ObsConfig.c1fb9e4f') : t('auto.ObsConfig.1e5a5df2')
)
const domainLabel = computed(() =>
  formData.provider === 'local' ? t('auto.ObsConfig.a78386b7') : 'CDN 域名（可选）'
)
const endpointPlaceholder = computed(() => {
  if (!formData.provider) return t('auto.ObsConfig.25f1b747')
  if (formData.provider === 'local') return './uploads'
  const map = {
    aliyun: 'https://xxx.oss-cn-hangzhou.aliyuncs.com',
    qiniu: 'http://xxx.bkt.clouddn.com',
    tencent: 'https://xxx.cos.ap-guangzhou.myqcloud.com',
    aws: 'https://s3.ap-southeast-1.amazonaws.com'
  }
  return map[formData.provider] || t('auto.ObsConfig.7d2b4ff6')
})
const domainPlaceholder = computed(() =>
  formData.provider === 'local' ? '/files' : t('auto.ObsConfig.a0218a65')
)

// 'default' 不在 el-tag 的合法 type 里（primary/success/info/warning/danger），会触发 Vue prop
// 校验告警并回落到 default='primary'。本地存储与未识别厂商用中性的 'info'。
const getProviderType = (provider) => {
  const types = { qiniu: 'warning', aliyun: 'primary', tencent: 'success', aws: 'info', local: 'info' }
  return types[provider] || 'info'
}
const getProviderLabel = (provider) => {
  const labels = { local: t('auto.ObsConfig.5e2bb608'), aliyun: t('auto.ObsConfig.a4bd7c03'), qiniu: t('auto.ObsConfig.3deda530'), tencent: t('auto.ObsConfig.b002b172'), aws: 'AWS S3' }
  return labels[provider] || provider
}

const handleProviderChange = (provider) => {
  formData.access_key = '';
  formData.secret_key = ''
  formData.bucket = ''
  formData.region = ''
  formData.endpoint = ''
  formData.domain = ''

  switch (provider) {
    case 'local':
      formData.endpoint = './uploads'
      formData.domain = '/files'
      break
    case 'aliyun':
      formData.endpoint = 'https://xxx.oss-cn-hangzhou.aliyuncs.com'
      break
    case 'qiniu':
      formData.endpoint = 'http://xxx.bkt.clouddn.com'
      break
    case 'tencent':
      formData.endpoint = 'https://xxx.cos.ap-guangzhou.myqcloud.com'
      break
    case 'aws':
      formData.region = 'ap-southeast-1'
      formData.endpoint = 'https://s3.ap-southeast-1.amazonaws.com'
      break
  }
};

const loadConfigList = async () => {
  loading.value = true
  try {
    const response = await getObsConfigList()
    configList.value = (response && response.list) || []
  } catch (error) {
    ElMessage.error(i18n.global.t(t('auto.ObsConfig.d98f57ba')))
  } finally {
    loading.value = false
  }
}

const handleCreate = () => {
  currentId.value = null
  resetForm()
  dialogVisible.value = true
}

const handleEdit = (row) => {
  currentId.value = row.id
  Object.assign(formData, row)
  dialogVisible.value = true
};

const handleTest = async (row) => {
  try {
    await testObsConnection(row.id)
    ElMessage.success(i18n.global.t(t('auto.ObsConfig.f20bc5e4')))
  } catch (error) {
    ElMessage.error(t('auto.ObsConfig.7b232863') + (error.message || t('auto.ObsConfig.974e7484')))
  }
}

const handleSetDefault = async (row) => {
  try {
    await setDefaultObsConfig(row.id)
    ElMessage.success(i18n.global.t(t('auto.ObsConfig.6355477c')))
    loadConfigList()
  } catch (error) {
    ElMessage.error(i18n.global.t(t('auto.ObsConfig.185d2e25')))
  }
}

const handleDelete = async (row) => {
  try {
    await ElMessageBox.confirm(t('auto.ObsConfig.7cfd11e6'), t('auto.ObsConfig.02d9819d'), { type: 'warning' })
    await deleteObsConfig(row.id)
    ElMessage.success(i18n.global.t(t('auto.ObsConfig.0007d170')))
    loadConfigList()
  } catch (error) {
    if (error !== 'cancel') ElMessage.error(i18n.global.t(t('auto.ObsConfig.acf0664a')))
  }
}

const handleSubmit = () => {
  formRef.value.validate(async (valid) => {
    if (!valid) return
    try {
      if (currentId.value) {
        await updateObsConfig(currentId.value, formData)
        ElMessage.success(i18n.global.t(t('auto.ObsConfig.55aa6366')))
      } else {
        await createObsConfig(formData)
        ElMessage.success(i18n.global.t(t('auto.ObsConfig.04a691b3')))
      }
      dialogVisible.value = false
      loadConfigList()
    } catch (error) {
      ElMessage.error(currentId.value ? t('auto.ObsConfig.930442e2') : t('auto.ObsConfig.a889286a'))
    }
  })
}

const handleDialogClose = () => {
  if (formRef.value) formRef.value.resetFields()
}

const resetForm = () => {
  Object.assign(formData, {
    name: '',
    provider: '',
    access_key: '',
    secret_key: '',
    bucket: '',
    region: '',
    endpoint: '',
    domain: '',
    path_prefix: '',
    max_size: 104857600,
    max_count: 1000,
    status: 'active'
  })
}

onMounted(() => {
  loadConfigList()
})
</script>

<style scoped>
.obs-config {
  padding: 20px;
}
.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}
.form-tip {
  margin-top: 5px;
  font-size: 12px;
  color: #909399;
  line-height: 1.5;
}
.muted {
  color: #c0c4cc;
}
</style>
