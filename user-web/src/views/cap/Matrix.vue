<template>
  <div class="app-container">
    <el-card shadow="never">
      <template #header>
        <div class="card-header">
          <span>{{ $t('auto.Matrix.6.fc94e6') }}</span>
          <el-button size="small" @click="copyMd">{{ $t('auto.Matrix.7.3006ed') }}</el-button>
        </div>
      </template>
      <el-alert
        type="info"
        :closable="false"
        show-icon
        :title="$t('auto.Matrix.1.e9242e')"
        style="margin-bottom: 12px"
      />
      <el-table v-loading="loading" :data="entries" border size="small">
        <el-table-column prop="platform" :label="$t('auto.Matrix.2.498128')" width="140" />
        <el-table-column prop="capability" :label="$t('auto.Matrix.3.89acdd')" width="160">
          <template #default="{ row }">
            <el-tag size="small">{{ row.capability }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="via" :label="$t('auto.Matrix.4.aaa159')" width="150" />
        <el-table-column prop="note" :label="$t('auto.Matrix.5.f411d0')" min-width="240" show-overflow-tooltip />
      </el-table>
    </el-card>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { ElMessage } from 'element-plus'
import { capabilitiesApi } from '@/api/kbConnector'

const { t } = useI18n()

const entries = ref([])
const markdown = ref('')
const loading = ref(false)

const load = async () => {
  loading.value = true
  try {
    const res = await capabilitiesApi.get()
    const data = res?.data ?? res ?? {}
    entries.value = data.entries ?? []
    markdown.value = data.markdown ?? ''
  } finally {
    loading.value = false
  }
}

const copyMd = async () => {
  try {
    await navigator.clipboard.writeText(markdown.value)
    ElMessage.success(t('auto.Matrix.8.52e6ab'))
  } catch {
    ElMessage.info(t('auto.Matrix.9.b2cb4b'))
  }
}

onMounted(load)
</script>

<style scoped>
.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}
</style>
