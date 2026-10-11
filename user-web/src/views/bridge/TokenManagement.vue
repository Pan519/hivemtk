<template>
  <div class="p-4">
    <el-card>
      <template #header>
        <div class="flex items-center justify-between">
          <span>{{ $t('auto.TokenManagement.5.ae69bb') }}</span>
          <el-tag :type="status.main_configured ? 'success' : 'danger'" size="small">
            {{ status.main_configured ? `已配置(${status.source === 'db' ? '运行时' : '环境变量'})` : '未配置(通道拒绝访问)' }}
          </el-tag>
        </div>
      </template>

      <el-alert v-if="status.prev_configured" type="warning" :closable="false" show-icon class="mb-3"
        :title="$t('auto.TokenManagement.1.736e03')" />

      <el-alert type="info" :closable="false" show-icon class="mb-3"
        :title="$t('auto.TokenManagement.2.d2a603')" />

      <el-form label-width="90px">
        <el-form-item :label="$t('auto.TokenManagement.3.8f9448')">
          <el-input v-model="newToken" readonly :placeholder="$t('auto.TokenManagement.4.78cfcb')">
            <template #append>
              <el-button @click="copyToken" :disabled="!newToken">{{ $t('auto.TokenManagement.6.79d3ab') }}</el-button>
            </template>
          </el-input>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" :loading="resetting" @click="onReset">{{ $t('auto.TokenManagement.7.456ef4') }}</el-button>
          <span class="ml-3 text-xs text-gray-400">{{ $t('auto.TokenManagement.8.fcf89e') }}</span>
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
    ElMessage.success(t('auto.TokenManagement.9.f2495b'))
    await loadStatus()
  } finally {
    resetting.value = false
  }
}

const copyToken = async () => {
  await navigator.clipboard.writeText(newToken.value)
  ElMessage.success(t('auto.TokenManagement.10.52e6ab'))
}

onMounted(loadStatus)
</script>
