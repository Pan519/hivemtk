<template>
  <div class="push-center">
    <el-row :gutter="16">
      <el-col :span="8" v-for="(status, plat) in quota" :key="plat">
        <el-card shadow="hover" class="quota-card">
          <div class="plat-name">{{ plat }}</div>
          <el-progress
            v-if="status.remain >= 0"
            :percentage="Math.max(0, 100 - (status.remain / (status.remain + 1) * 100))"
            :status="status.remain > 100 ? 'success' : status.remain > 20 ? 'warning' : 'exception'"
          />
          <div class="plat-detail">
            <el-tag v-if="status.registered" type="success">{{ $t('auto.PushCenter.19f13a7c') }}</el-tag>
            <el-tag v-else type="info">{{ $t('auto.PushCenter.0e6fe5bc') }}</el-tag>
            <span v-if="status.remain >= 0">剩余 {{ status.remain }}</span>
            <span v-else>{{ $t('auto.PushCenter.bac375d8') }}</span>
          </div>
        </el-card>
      </el-col>
    </el-row>

    <el-card shadow="never" style="margin-top:16px">
      <template #header>
        <div class="card-header">
          <span>{{ $t('auto.PushCenter.2fcaefca') }}</span>
        </div>
      </template>
      <el-input
        v-model="urlsText"
        type="textarea"
        :rows="4"
        :placeholder="$t('auto.PushCenter.9af24eca')"
      />
      <div style="margin-top:8px">
        <el-checkbox-group v-model="selectedPlatforms">
          <el-checkbox value="baidu">{{ $t('auto.PushCenter.a3f4a5b0') }}</el-checkbox>
          <el-checkbox value="indexnow">IndexNow (Bing/Yandex/Naver)</el-checkbox>
          <el-checkbox value="google">Google</el-checkbox>
        </el-checkbox-group>
        <el-button type="primary" @click="push" style="margin-left:12px">{{ $t('auto.PushCenter.48aa9bc8') }}</el-button>
      </div>
    </el-card>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { geoApi } from '@/api/geo'
import { ElMessage } from 'element-plus'
const { t } = useI18n()

const quota = ref({})
const urlsText = ref('')
const selectedPlatforms = ref(['indexnow', 'baidu'])
async function loadQuota() {
  try {
    quota.value = (await geoApi.getPushQuota()) || {}
  } catch (e) { /* 忽略异常：失败时保持既有状态，不打断用户 */ }
}
async function push() {
  const urls = urlsText.value.split('\n').map(s => s.trim()).filter(Boolean)
  if (!urls.length) return
  try {
    await geoApi.pushUrls({ urls, platforms: selectedPlatforms.value })
    ElMessage.success(t('auto.PushCenter.3aefb579'))
    urlsText.value = ''
    loadQuota()
  } catch (e) { ElMessage.error(t('auto.PushCenter.3b6ba072')) }
}
onMounted(loadQuota)
</script>

<style scoped>
.quota-card { text-align:center; }
.plat-name { font-size:16px; font-weight:bold; margin-bottom:8px; }
.plat-detail { margin-top:8px; font-size:12px; color:#666; }
</style>
