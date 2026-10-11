<template>
  <div class="seo-infra">
    <el-tabs v-model="tab">
      <el-tab-pane label="llms.txt v2" name="llms">
        <div style="display:flex;gap:16px">
          <div style="flex:1">
            <h4>{{ $t('auto.SeoInfra.1.645dbc') }}</h4>
            <pre class="code">{{ llmsPreview }}</pre>
          </div>
          <div style="flex:1">
            <h4>{{ $t('auto.SeoInfra.2.4ffba8') }}</h4>
            <ul>
              <li>{{ $t('auto.SeoInfra.3.50b864') }}</li>
              <li>{{ $t('auto.SeoInfra.4.0aec73') }}</li>
              <li>{{ $t('auto.SeoInfra.5.919198') }}</li>
              <li>{{ $t('auto.SeoInfra.6.2a4fe4') }}</li>
              <li>{{ $t('auto.SeoInfra.7.74d369') }}</li>
              <li>{{ $t('auto.SeoInfra.8.8592ec') }}</li>
            </ul>
          </div>
        </div>
      </el-tab-pane>
      <el-tab-pane label="robots.txt" name="robots">
        <pre class="code">{{ robotsPreview }}</pre>
      </el-tab-pane>
      <el-tab-pane label="sitemap.xml" name="sitemap">
        <el-button @click="genSitemap" style="margin-bottom:8px">{{ $t('auto.SeoInfra.9.4dfe70') }}</el-button>
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
