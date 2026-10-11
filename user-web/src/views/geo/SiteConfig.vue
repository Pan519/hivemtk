<template>
  <div class="geo-site-config">
    <el-card shadow="never">
      <template #header>
        <div class="card-header">
          <span>{{ $t('auto.SiteConfig.4d92789a') }}</span>
          <el-button type="primary" @click="openDialog()">{{ $t('auto.SiteConfig.8604ec54') }}</el-button>
        </div>
      </template>

      <el-table :data="list" border stripe size="small">
        <el-table-column prop="domain" :label="$t('auto.SiteConfig.190980c7')" />
        <el-table-column prop="provider" :label="$t('auto.SiteConfig.8da5f75a')" width="120">
          <template #default="{ row }">
            <el-tag size="small">{{ row.provider || 'cloudflare' }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="hugo_path" :label="$t('auto.SiteConfig.b078ffb4')" show-overflow-tooltip />
        <el-table-column prop="git_repo" :label="$t('auto.SiteConfig.d71d035f')" show-overflow-tooltip />
        <el-table-column :label="$t('auto.SiteConfig.fbc026c9')" width="80">
          <template #default="{ row }">
            <el-tag :type="row.health_score > 80 ? 'success' : row.health_score > 40 ? 'warning' : 'danger'" size="small">
              {{ row.health_score || 0 }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column :label="$t('auto.SiteConfig.3fea7ca7')" width="80">
          <template #default="{ row }">
            <el-tag :type="row.active ? 'success' : 'info'" size="small">
              {{ row.active ? '启用' : '停用' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column :label="$t('auto.SiteConfig.2b6bc0f2')" width="180" fixed="right">
          <template #default="{ row }">
            <el-button size="small" @click="openDialog(row)">{{ $t('auto.SiteConfig.95b351c8') }}</el-button>
            <el-button size="small" type="danger" @click="remove(row)">{{ $t('auto.SiteConfig.2f4aaddd') }}</el-button>
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
    <el-dialog v-model="dialogVisible" ::title="$t('auto.SiteConfig.4275a15b')" width="560px">
      <el-form :model="form" label-width="100px">
        <el-form-item :label="$t('auto.SiteConfig.190980c7')" required>
          <el-input v-model="form.domain" placeholder="example.com" />
        </el-form-item>
        <el-form-item :label="$t('auto.SiteConfig.8da5f75a')">
          <el-select v-model="form.provider" style="width:100%">
            <el-option label="Cloudflare Pages" value="cloudflare" />
            <el-option label="GitHub Pages" value="github" />
            <el-option label="Vercel" value="vercel" />
            <el-option :label="$t('auto.SiteConfig.04d686ca')" value="self" />
          </el-select>
        </el-form-item>
        <el-form-item :label="$t('auto.SiteConfig.b078ffb4')">
          <el-input v-model="form.hugo_path" placeholder="/data/hivemtk-site" />
        </el-form-item>
        <el-form-item :label="$t('auto.SiteConfig.d71d035f')">
          <el-input v-model="form.git_repo" placeholder="git@github.com:xxx/xxx.git" />
        </el-form-item>
        <el-form-item :label="$t('auto.SiteConfig.65b20ba7')">
          <el-input v-model="form.deploy_cmd" placeholder="hugo -e production" />
        </el-form-item>
        <el-form-item :label="$t('auto.SiteConfig.7854b52a')">
          <el-switch v-model="form.active" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">{{ $t('auto.SiteConfig.625fb26b') }}</el-button>
        <el-button type="primary" @click="save">{{ $t('auto.SiteConfig.be5fbbe3') }}</el-button>
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
  id: '', domain: '', provider: 'cloudflare',
  hugo_path: '', git_repo: '', deploy_cmd: '', active: true,
})

function resetForm() {
  form.value = {
    id: '', domain: '', provider: 'cloudflare',
    hugo_path: '', git_repo: '', deploy_cmd: '', active: true,
  }
  isEdit.value = false
}

async function load() {
  try {
    const res = await geoApi.listSites({ page: page.value, limit: limit.value })
    list.value = res?.list || []
    total.value = res?.total || 0
  } catch (e) { /* 忽略异常：失败时保持既有状态，不打断用户 */ }
}

function openDialog(row) {
  if (row) {
    form.value = { ...row }
    isEdit.value = true
  } else {
    resetForm()
  }
  dialogVisible.value = true
}

async function save() {
  try {
    if (isEdit.value) {
      await geoApi.updateSite(form.value.id, form.value)
    } else {
      await geoApi.createSite(form.value)
    }
    ElMessage.success(t('auto.SiteConfig.f8dfedcd'))
    dialogVisible.value = false
    load()
  } catch (e) {
    ElMessage.error(e?.response?.data?.message || t('auto.SiteConfig.6de920b4'))
  }
}

async function remove(row) {
  try {
    await ElMessageBox.confirm(`确认删除 ${row.domain} ?`, t('auto.SiteConfig.900c70fa'), { type: 'warning' })
    await geoApi.deleteSite(row.id)
    ElMessage.success(t('auto.SiteConfig.5cc23262'))
    load()
  } catch (e) { /* 忽略异常：失败时保持既有状态，不打断用户 */ }
}

onMounted(load)
</script>

<style scoped>
.card-header { display: flex; justify-content: space-between; align-items: center; }
.pager { display: flex; justify-content: flex-end; margin-top: 12px; }
</style>
