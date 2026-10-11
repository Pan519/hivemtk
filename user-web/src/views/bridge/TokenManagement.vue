<template>
  <div class="p-4">
    <el-card>
      <template #header>
        <div class="flex items-center justify-between">
          <span>{{ $t('auto.TokenManagement.ae69bbab') }}</span>
          <el-tag :type="status.main_configured ? 'success' : 'danger'" size="small">
            {{ status.main_configured ? `已配置(${status.source === 'db' ? '运行时' : '环境变量'})` : '未配置(通道拒绝访问)' }}
          </el-tag>
        </div>
      </template>

      <el-alert v-if="status.prev_configured" type="warning" :closable="false" show-icon class="mb-3"
        :title="$t('auto.TokenManagement.736e03bb')" />

      <el-alert type="info" :closable="false" show-icon class="mb-3"
        :title="$t('auto.TokenManagement.d2a60388')" />

      <el-form label-width="90px">
        <el-form-item :label="$t('auto.TokenManagement.8f944893')">
          <el-input v-model="newToken" readonly :placeholder="$t('auto.TokenManagement.78cfcb87')">
            <template #append>
              <el-button @click="copyToken" :disabled="!newToken">{{ $t('auto.TokenManagement.79d3abe9') }}</el-button>
            </template>
          </el-input>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" :loading="resetting" @click="onReset">{{ $t('auto.TokenManagement.456ef430') }}</el-button>
          <span class="ml-3 text-xs text-gray-400">{{ $t('auto.TokenManagement.fcf89e1b') }}</span>
        </el-form-item>
      </el-form>
    </el-card>
  </div>
</template>

<script setup>
import { onMounted, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { ElMessage } from 'element-plus'
import { getBridgeTokenStatus, resetBridgeToken } from '@/api/bridgeToken.js'

const { t } = useI18n()

const status = ref({})
const newToken = ref('')
const resetting = ref(false)

const loadStatus = async () => {
  try {
    status.value = (await getBridgeTokenStatus()) || {}
  } catch (e) { console.warn("[request] 后台调用失败(已忽略):", e) }
}

const onReset = async () => {
  resetting.value = true
  try {
    const res = await resetBridgeToken()
    newToken.value = res?.token || ''
    ElMessage.success(t('auto.TokenManagement.f2495b61'))
    await loadStatus()
  } finally {
    resetting.value = false
  }
}

const copyToken = async () => {
  await navigator.clipboard.writeText(newToken.value)
  ElMessage.success(t('auto.TokenManagement.52e6abbe'))
}

onMounted(loadStatus)
</script>
