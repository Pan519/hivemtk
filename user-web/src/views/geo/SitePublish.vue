<template>
  <div class="site-publish">
    <el-card shadow="never">
      <template #header>
        <div class="card-header">
          <span>{{ $t('auto.SitePublish.5.53ae7c') }}</span>
        </div>
      </template>
      <el-steps :active="step" finish-status="success" simple style="margin-bottom:24px">
        <el-step :title="$t('auto.SitePublish.1.1d6248')" />
        <el-step :title="$t('auto.SitePublish.2.f38989')" />
        <el-step :title="$t('auto.SitePublish.3.e90478')" />
        <el-step :title="$t('auto.SitePublish.4.13b73c')" />
      </el-steps>

      <el-button-group style="margin-bottom:16px">
        <el-button @click="exportOnly">{{ $t('auto.SitePublish.6.48885d') }}</el-button>
        <el-button @click="deployOnly">{{ $t('auto.SitePublish.7.e500b7') }}</el-button>
        <el-button type="primary" @click="fullPipeline">{{ $t('auto.SitePublish.8.cd7088') }}</el-button>
      </el-button-group>

      <el-alert v-if="lastResult" :title="lastResult" type="success" show-icon :closable="false" />
    </el-card>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { geoApi } from '@/api/geo'
import { ElMessage } from 'element-plus'
const { t } = useI18n()

const step = ref(0)
const lastResult = ref('')
async function exportOnly() {
  try {
    const res = await geoApi.exportSite()
    step.value = 1
    lastResult.value = `导出 ${res?.exported ?? 0} 篇`
    ElMessage.success(t('auto.SitePublish.9.105c8a'))
  } catch (e) { ElMessage.error(t('auto.SitePublish.10.dd51ab')) }
}
async function deployOnly() {
  try {
    const res = await geoApi.deploySite()
    step.value = 2
    lastResult.value = `部署地址：${res?.site_url || '—'}`
    ElMessage.success(t('auto.SitePublish.11.446c6b'))
  } catch (e) { ElMessage.error(t('auto.SitePublish.12.1cdfe9')) }
}
async function fullPipeline() {
  step.value = 1
  try {
    const res = await geoApi.siteFullPipeline()
    step.value = 4
    lastResult.value = `导出 ${res?.exported} 篇 → ${res?.site_url}`
    ElMessage.success(t('auto.SitePublish.13.aafd24'))
  } catch (e) { step.value = 0 }
}
</script>
