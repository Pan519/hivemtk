<template>
  <div class="geo-page p-4">
    <el-card>
      <template #header>
        <div class="flex items-center justify-between">
          <span class="font-bold">{{ $t('auto.CompetitorManage.15.83cab0') }}</span>
          <el-button type="primary" size="small" @click="openCreate">{{ $t('auto.CompetitorManage.16.5da4d6') }}</el-button>
        </div>
      </template>

      <el-table :data="list" v-loading="loading" size="default" stripe>
        <el-table-column prop="id" label="ID" width="60" />
        <el-table-column prop="name" :label="$t('auto.CompetitorManage.1.d7ec2d')" width="140" />
        <el-table-column prop="domain" :label="$t('auto.CompetitorManage.2.190980')" width="220">
          <template #default="{ row }">
            <a :href="'https://' + row.domain" target="_blank" class="font-mono text-blue-600">{{ row.domain }}</a>
          </template>
        </el-table-column>
        <el-table-column :label="$t('auto.CompetitorManage.3.6147cf')" min-width="240">
          <template #default="{ row }">
            <div class="flex flex-wrap gap-1">
              <el-tag v-for="p in (row.paths || [])" :key="p" size="small" effect="plain">{{ p }}</el-tag>
            </div>
          </template>
        </el-table-column>
        <el-table-column prop="category" :label="$t('auto.CompetitorManage.4.d0771a')" width="100">
          <template #default="{ row }">
            <el-tag size="small" :type="categoryTag(row.category)">{{ categoryLabel(row.category) }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="priority" :label="$t('auto.CompetitorManage.5.ee8ecb')" width="80" sortable />
        <el-table-column prop="status" :label="$t('auto.CompetitorManage.6.3fea7c')" width="90">
          <template #default="{ row }">
            <el-tag size="small" :type="row.status === 'active' ? 'success' : 'info'">{{ row.status }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column :label="$t('auto.CompetitorManage.7.2b6bc0')" width="180">
          <template #default="{ row }">
            <el-button size="small" @click="openEdit(row)">{{ $t('auto.CompetitorManage.17.95b351') }}</el-button>
            <el-button size="small" type="danger" plain @click="onDelete(row)">{{ $t('auto.CompetitorManage.18.2f4aad') }}</el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    
    <el-dialog v-model="dialogVisible" ::title="$t('auto.CompetitorManage.8.e2045d')" width="520px" destroy-on-close>
      <el-form :model="form" label-width="90px" size="default">
        <el-form-item :label="$t('auto.CompetitorManage.1.d7ec2d')">
          <el-input v-model="form.name" :placeholder="$t('auto.CompetitorManage.9.68edbd')" />
        </el-form-item>
        <el-form-item :label="$t('auto.CompetitorManage.2.190980')">
          <el-input v-model="form.domain" :placeholder="$t('auto.CompetitorManage.10.0aa1c6')" :disabled="!!editing" />
        </el-form-item>
        <el-form-item :label="$t('auto.CompetitorManage.3.6147cf')">
          <div>
            <el-tag v-for="(p, i) in form.paths" :key="i" closable @close="form.paths.splice(i, 1)" size="default" style="margin-right:4px;margin-bottom:4px">
              {{ p }}
            </el-tag>
            <el-input v-if="newPath" v-model="newPath" size="small" style="width:120px" @keyup.enter="addPath" />
            <el-button v-else size="small" plain @click="newPath = '/'">{{ $t('auto.CompetitorManage.19.047037') }}</el-button>
          </div>
        </el-form-item>
        <el-form-item :label="$t('auto.CompetitorManage.4.d0771a')">
          <el-select v-model="form.category">
            <el-option :label="$t('auto.CompetitorManage.11.710d58')" value="direct" />
            <el-option :label="$t('auto.CompetitorManage.12.80203f')" value="global" />
            <el-option :label="$t('auto.CompetitorManage.13.307609')" value="indirect" />
          </el-select>
        </el-form-item>
        <el-form-item :label="$t('auto.CompetitorManage.5.ee8ecb')">
          <el-slider v-model="form.priority" :min="1" :max="10" show-input style="width:240px" />
        </el-form-item>
        <el-form-item :label="$t('auto.CompetitorManage.6.3fea7c')">
          <el-switch v-model="form.status" active-value="active" inactive-value="paused" active-text="active" inactive-text="paused" />
        </el-form-item>
        <el-form-item :label="$t('auto.CompetitorManage.14.2432b5')">
          <el-input v-model="form.notes" type="textarea" :rows="2" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">{{ $t('auto.CompetitorManage.20.625fb2') }}</el-button>
        <el-button type="primary" :loading="saving" @click="onSave">{{ $t('auto.CompetitorManage.21.be5fbb') }}</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { ElMessage, ElMessageBox } from 'element-plus'
import { listCompetitors, createCompetitor, updateCompetitor, deleteCompetitor } from '@/api/geoProbe.js'

const { t } = useI18n()

const list = ref([])
const loading = ref(false)
const dialogVisible = ref(false)
const editing = ref(null)
const saving = ref(false)
const newPath = ref('')

const form = ref({
  name: '', domain: '', paths: ['/', '/product'],
  category: 'direct', priority: 5, status: 'active', notes: ''
})

const categoryLabel = (c) => ({ direct: t('auto.CompetitorManage.11.710d58'), global: t('auto.CompetitorManage.12.80203f'), indirect: t('auto.CompetitorManage.13.307609') }[c] || c || '—')
const categoryTag = (c) => ({ direct: 'success', global: 'warning', indirect: 'info' }[c] || 'info')

async function load() {
  loading.value = true
  try {
    list.value = await listCompetitors()
  } finally {
    loading.value = false
  }
}

function resetForm() {
  form.value = {
    name: '', domain: '', paths: ['/', '/product'],
    category: 'direct', priority: 5, status: 'active', notes: ''
  }
  newPath.value = ''
}

function openCreate() {
  editing.value = null
  resetForm()
  dialogVisible.value = true
}

function openEdit(row) {
  editing.value = row
  form.value = JSON.parse(JSON.stringify(row))
  if (!form.value.paths || form.value.paths.length === 0) {
    form.value.paths = ['/']
  }
  dialogVisible.value = true
}

function addPath() {
  if (newPath.value && !form.value.paths.includes(newPath.value)) {
    form.value.paths.push(newPath.value)
  }
  newPath.value = ''
}

async function onSave() {
  if (!form.value.name || !form.value.domain) {
    ElMessage.warning(t('auto.CompetitorManage.22.59c491'))
    return
  }
  saving.value = true
  try {
    if (editing.value) {
      await updateCompetitor(editing.value.id, form.value)
      ElMessage.success(t('auto.CompetitorManage.23.55aa63'))
    } else {
      await createCompetitor(form.value)
      ElMessage.success(t('auto.CompetitorManage.24.04a691'))
    }
    dialogVisible.value = false
    await load()
  } catch (e) {
    ElMessage.error(e?.response?.data?.message || t('auto.CompetitorManage.25.6de920'))
  } finally {
    saving.value = false
  }
}

async function onDelete(row) {
  try {
    await ElMessageBox.confirm(`确认删除竞品 "${row.name}" 吗？`, t('auto.CompetitorManage.26.02d981'), { type: 'warning' })
    await deleteCompetitor(row.id)
    ElMessage.success(t('auto.CompetitorManage.27.5cc232'))
    await load()
  } catch { /* 忽略：清理/存储/恢复类 best-effort 操作 */ }
}

onMounted(load)
</script>
