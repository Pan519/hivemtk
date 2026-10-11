<template>
  <div class="geo-pusher-config">
    <el-card shadow="never">
      <template #header>
        <div class="card-header">
          <span>{{ $t('auto.PusherConfig.13.e417fe') }}</span>
          <el-button type="primary" @click="openDialog()">{{ $t('auto.PusherConfig.14.4fd331') }}</el-button>
        </div>
      </template>

      <el-alert type="info" :closable="false" style="margin-bottom:12px">
        <template #title>
          凭据用 AES-256-GCM 加密后存 JSONB。百度：token + site；Google：多个 Service Account；IndexNow：host + key（无限额）
        </template>
      </el-alert>

      <el-table :data="list" border stripe size="small">
        <el-table-column prop="platform" :label="$t('auto.PusherConfig.1.6c1b19')" width="120">
          <template #default="{ row }">
            <el-tag size="small" :type="platformTag(row.platform)">{{ row.platform }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column :label="$t('auto.PusherConfig.2.f74c1d')" width="160">
          <template #default="{ row }">
            <el-progress
              :percentage="row.daily_limit > 0 ? Math.min(100, Math.round((row.used_today / row.daily_limit) * 100)) : 0"
              :status="row.daily_limit - row.used_today > 500 ? 'success' : row.daily_limit - row.used_today > 100 ? 'warning' : 'exception'"
              :format="() => row.daily_limit > 0 ? `${row.used_today}/${row.daily_limit}` : '∞'"
            />
          </template>
        </el-table-column>
        <el-table-column :label="$t('auto.PusherConfig.3.3fea7c')" width="100">
          <template #default="{ row }">
            <el-tag :type="row.active ? 'success' : 'info'" size="small">{{ row.active ? '启用' : '停用' }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="updated_at" :label="$t('auto.PusherConfig.4.a001a2')" width="180" />
        <el-table-column :label="$t('auto.PusherConfig.5.2b6bc0')" width="180" fixed="right">
          <template #default="{ row }">
            <el-button size="small" @click="openDialog(row)">{{ $t('auto.PusherConfig.15.95b351') }}</el-button>
            <el-button size="small" type="danger" @click="remove(row)">{{ $t('auto.PusherConfig.16.2f4aad') }}</el-button>
          </template>
        </el-table-column>
      </el-table>

      <div class="pager">
        <el-pagination
          v-model:current-page="page"
          v-model:page-size="limit"
          :total="total"
          :page-sizes="[20, 50, 100]"
          layout="total, sizes, prev, pager, next"
          @size-change="load"
          @current-change="load"
        />
      </div>
    </el-card>

    <!-- Dialog -->
    <el-dialog v-model="dialogVisible" ::title="$t('auto.PusherConfig.6.13e924')" width="560px">
      <el-form :model="form" label-width="100px">
        <el-form-item :label="$t('auto.PusherConfig.7.498128')" required>
          <el-select v-model="form.platform" :disabled="isEdit" style="width:100%">
            <el-option :label="$t('auto.PusherConfig.8.32cc42')" value="baidu" />
            <el-option label="Google (Indexing API)" value="google" />
            <el-option label="IndexNow (Bing/Yandex/Naver)" value="indexnow" />
            <el-option :label="$t('auto.PusherConfig.9.357b88')" value="toutiao" />
            <el-option :label="$t('auto.PusherConfig.10.852a9b')" value="shenma" />
          </el-select>
        </el-form-item>
        <el-form-item :label="$t('auto.PusherConfig.11.dc861f')">
          <el-input
            v-model="form.config_json_text"
            type="textarea"
            :rows="5"
            placeholder='{"site":"example.com","token":"xxx"}'
          />
        </el-form-item>
        <el-form-item :label="$t('auto.PusherConfig.2.f74c1d')">
          <el-input-number v-model="form.daily_limit" :min="0" :max="100000" />
          <span style="margin-left:8px;color:#999">{{ $t('auto.PusherConfig.17.125aaa') }}</span>
        </el-form-item>
        <el-form-item :label="$t('auto.PusherConfig.12.7854b5')">
          <el-switch v-model="form.active" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">{{ $t('auto.PusherConfig.18.625fb2') }}</el-button>
        <el-button type="primary" @click="save">{{ $t('auto.PusherConfig.19.be5fbb') }}</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { geoApi } from '@/api/geo'
import { ElMessage, ElMessageBox } from 'element-plus'

const { t } = useI18n()

const list = ref([])
const page = ref(1)
const limit = ref(20)
const total = ref(0)
const dialogVisible = ref(false)
const isEdit = ref(false)
const form = ref({
  id: '', platform: '', config_json_text: '', daily_limit: 0, active: true,
})

function resetForm() {
  form.value = { id: '', platform: '', config_json_text: '', daily_limit: 0, active: true }
  isEdit.value = false
}

function platformTag(p) {
  const m = { baidu: 'warning', google: 'primary', indexnow: 'success', toutiao: 'info', shenma: 'info' }
  return m[p] || ''
}

async function load() {
  try {
    const res = await geoApi.listPushers({ page: page.value, limit: limit.value })
    list.value = res?.list || []
    total.value = res?.total || 0
  } catch (e) { /* 忽略异常：失败时保持既有状态，不打断用户 */ }
}

function openDialog(row) {
  if (row) {
    form.value = {
      id: row.id, platform: row.platform, daily_limit: row.daily_limit, active: row.active,
      config_json_text: typeof row.config_json === 'string' ? row.config_json : JSON.stringify(row.config_json || {}),
    }
    isEdit.value = true
  } else {
    resetForm()
  }
  dialogVisible.value = true
}

async function save() {
  try {
    // 把 config_json_text 解析成对象再序列化
    let cfgJson = form.value.config_json_text
    try { cfgJson = JSON.parse(form.value.config_json_text) } catch (e) { /* 忽略异常：失败时保持既有状态，不打断用户 */ }
    const payload = {
      platform: form.value.platform,
      config_json: cfgJson,
      daily_limit: form.value.daily_limit,
      active: form.value.active,
    }
    if (isEdit.value) {
      await geoApi.updatePusher(form.value.id, payload)
    } else {
      await geoApi.createPusher(payload)
    }
    ElMessage.success(t('auto.PusherConfig.20.f8dfed'))
    dialogVisible.value = false
    load()
  } catch (e) {
    ElMessage.error(e?.response?.data?.message || t('auto.PusherConfig.21.6de920'))
  }
}

async function remove(row) {
  try {
    await ElMessageBox.confirm(`确认删除 ${row.platform} ?`, t('auto.PusherConfig.22.900c70'), { type: 'warning' })
    await geoApi.deletePusher(row.id)
    ElMessage.success(t('auto.PusherConfig.23.5cc232'))
    load()
  } catch (e) { /* 忽略异常：失败时保持既有状态，不打断用户 */ }
}

onMounted(load)
</script>

<style scoped>
.card-header { display: flex; justify-content: space-between; align-items: center; }
.pager { display: flex; justify-content: flex-end; margin-top: 12px; }
</style>
