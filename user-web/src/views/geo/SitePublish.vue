<template>
  <div class="site-publish">
    <el-card shadow="never">
      <template #header>
        <div class="card-header">
          <span>{{ $t('auto.SitePublish.53ae7c9e') }}</span>
        </div>
      </template>
      <el-steps :active="step" finish-status="success" simple style="margin-bottom:24px">
        <el-step :title="$t('auto.SitePublish.1d624848')" />
        <el-step :title="$t('auto.SitePublish.f3898916')" />
        <el-step :title="$t('auto.SitePublish.e90478d7')" />
        <el-step :title="$t('auto.SitePublish.13b73c87')" />
      </el-steps>

      <el-button-group style="margin-bottom:16px">
        <el-button @click="exportOnly">{{ $t('auto.SitePublish.48885dbf') }}</el-button>
        <el-button @click="deployOnly">{{ $t('auto.SitePublish.e500b7da') }}</el-button>
        <el-button type="primary" @click="fullPipeline">{{ $t('auto.SitePublish.cd7088c2') }}</el-button>
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
    ElMessage.success(t('auto.SitePublish.105c8a51'))
  } catch (e) { ElMessage.error(t('auto.SitePublish.dd51ab50')) }
}
async function deployOnly() {
  try {
    const res = await geoApi.deploySite()
    step.value = 2
    lastResult.value = `部署地址：${res?.site_url || '—'}`
    ElMessage.success(t('auto.SitePublish.446c6b6d'))
  } catch (e) { ElMessage.error(t('auto.SitePublish.1cdfe97c')) }
}
async function fullPipeline() {
  step.value = 1
  try {
    const res = await geoApi.siteFullPipeline()
    step.value = 4
    lastResult.value = `导出 ${res?.exported} 篇 → ${res?.site_url}`
    ElMessage.success(t('auto.SitePublish.aafd2400'))
  } catch (e) { step.value = 0 }
}
</script>
