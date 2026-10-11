<template>
  <div class="geo-brand-config">
    <el-tabs v-model="activeTab">
      <el-tab-pane :label="$t('auto.BrandConfig.1.05b4f3')" name="llms">
        <el-card shadow="never">
          <template #header>
            <div class="card-header">
              <span>{{ $t('auto.BrandConfig.5.ff9c4c') }}</span>
              <div>
                <el-input v-model="previewDomain" placeholder="example.com" style="width:200px;margin-right:8px" />
                <el-button type="primary" @click="loadLlms">{{ $t('auto.BrandConfig.6.b542b0') }}</el-button>
              </div>
            </div>
          </template>
          <el-empty v-if="!llmsContent" :description="$t('auto.BrandConfig.2.94ef7c')">
            <el-button type="primary" @click="loadLlms">{{ $t('auto.BrandConfig.7.50e17c') }}</el-button>
          </el-empty>
          <pre v-else class="code-block">{{ llmsContent }}</pre>
          <div v-if="llmsContent" style="margin-top:12px">
            <el-tag type="success">{{ $t('auto.BrandConfig.8.38bc6a') }}</el-tag>
            <el-tag style="margin-left:8px">{{ $t('auto.BrandConfig.9.5ef0f0') }}</el-tag>
            <el-tag style="margin-left:8px">{{ $t('auto.BrandConfig.10.333724') }}</el-tag>
          </div>
        </el-card>
      </el-tab-pane>

      <el-tab-pane :label="$t('auto.BrandConfig.3.c91198')" name="robots">
        <el-card shadow="never">
          <template #header>
            <div class="card-header">
              <span>{{ $t('auto.BrandConfig.11.a14124') }}</span>
              <el-button type="primary" @click="loadRobots">{{ $t('auto.BrandConfig.6.b542b0') }}</el-button>
            </div>
          </template>
          <el-empty v-if="!robotsContent" :description="$t('auto.BrandConfig.4.1c0c7e')" />
          <pre v-else class="code-block">{{ robotsContent }}</pre>
          <div v-if="robotsContent" style="margin-top:12px">
            <el-tag type="success">{{ $t('auto.BrandConfig.12.d7460e') }}</el-tag>
            <el-tag style="margin-left:8px">Bytespider/DeepSeekBot/QwenBot/ErnieBot</el-tag>
            <el-tag style="margin-left:8px">GPTBot/ClaudeBot/PerplexityBot/Amazonbot</el-tag>
          </div>
        </el-card>
      </el-tab-pane>
    </el-tabs>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { geoApi } from '@/api/geo'

const { t } = useI18n()

const activeTab = ref('llms')
const previewDomain = ref('example.com')
const llmsContent = ref('')
const robotsContent = ref('')

async function loadLlms() {
  try {
    const res = await geoApi.getLlmsTxtPreview(previewDomain.value)
    llmsContent.value = res?.content || ''
  } catch (e) {
    llmsContent.value = '（后端暂未部署 llms.txt 生成逻辑）'
  }
}

async function loadRobots() {
  try {
    const res = await geoApi.getRobotsTxtPreview()
    robotsContent.value = res?.content || ''
  } catch (e) {
    robotsContent.value = '（后端暂未部署 robots.txt 生成逻辑）'
  }
}

onMounted(() => {
  loadLlms()
  loadRobots()
})
</script>

<style scoped>
.card-header { display: flex; justify-content: space-between; align-items: center; }
.code-block {
  background: #1e1e2e;
  color: #cdd6f4;
  padding: 16px;
  border-radius: 8px;
  font-family: 'JetBrains Mono', monospace;
  font-size: 12px;
  white-space: pre-wrap;
  word-break: break-all;
  max-height: 400px;
  overflow-y: auto;
}
</style>
