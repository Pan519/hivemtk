<template>
  <div class="seo-infra">
    <el-tabs v-model="tab">
      <el-tab-pane label="llms.txt v2" name="llms">
        <div style="display:flex;gap:16px">
          <div style="flex:1">
            <h4>{{ $t('auto.SeoInfra.645dbc55') }}</h4>
            <pre class="code">{{ llmsPreview }}</pre>
          </div>
          <div style="flex:1">
            <h4>{{ $t('auto.SeoInfra.4ffba880') }}</h4>
            <ul>
              <li>{{ $t('auto.SeoInfra.50b8640f') }}</li>
              <li>{{ $t('auto.SeoInfra.0aec7373') }}</li>
              <li>{{ $t('auto.SeoInfra.9191988a') }}</li>
              <li>{{ $t('auto.SeoInfra.2a4fe4bc') }}</li>
              <li>{{ $t('auto.SeoInfra.74d36944') }}</li>
              <li>{{ $t('auto.SeoInfra.8592ec89') }}</li>
            </ul>
          </div>
        </div>
      </el-tab-pane>
      <el-tab-pane label="robots.txt" name="robots">
        <pre class="code">{{ robotsPreview }}</pre>
      </el-tab-pane>
      <el-tab-pane label="sitemap.xml" name="sitemap">
        <el-button @click="genSitemap" style="margin-bottom:8px">{{ $t('auto.SeoInfra.4dfe7036') }}</el-button>
        <pre class="code">{{ sitemapPreview }}</pre>
      </el-tab-pane>
    </el-tabs>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { geoApi } from '@/api/geo'
const { t } = useI18n()

const tab = ref('llms')
const llmsPreview = ref('')
const robotsPreview = ref('')
const sitemapPreview = ref('')
async function load() {
  try {
    const [l, r] = await Promise.all([
      geoApi.getLlmsTxtPreview(),
      geoApi.getRobotsTxtPreview(),
    ])
    llmsPreview.value = l?.content || ''
    robotsPreview.value = r?.content || ''
  } catch (e) { /* 忽略异常：失败时保持既有状态，不打断用户 */ }
}
async function genSitemap() {
  try {
    const r = await geoApi.getSitemapPreview()
    sitemapPreview.value = r?.content || ''
  } catch (e) { /* 忽略异常：失败时保持既有状态，不打断用户 */ }
}
onMounted(load)
</script>

<style scoped>
.code { background:#1e1e1e; color:#d4d4d4; padding:16px; border-radius:6px; max-height:400px; overflow:auto; font-size:12px; }
</style>
