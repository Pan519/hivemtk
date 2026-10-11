<template>
  <div class="geo-page p-4">
    <el-card>
      <template #header>
        <div class="flex items-center justify-between">
          <span class="font-bold">{{ $t('auto.CompetitorManage.83cab0e2') }}</span>
          <el-button type="primary" size="small" @click="openCreate">{{ $t('auto.CompetitorManage.5da4d69d') }}</el-button>
        </div>
      </template>

      <el-table :data="list" v-loading="loading" size="default" stripe>
        <el-table-column prop="id" label="ID" width="60" />
        <el-table-column prop="name" :label="$t('auto.CompetitorManage.d7ec2d3f')" width="140" />
        <el-table-column prop="domain" :label="$t('auto.CompetitorManage.190980c7')" width="220">
          <template #default="{ row }">
            <a :href="'https://' + row.domain" target="_blank" class="font-mono text-blue-600">{{ row.domain }}</a>
          </template>
        </el-table-column>
        <el-table-column :label="$t('auto.CompetitorManage.6147cfcd')" min-width="240">
          <template #default="{ row }">
            <div class="flex flex-wrap gap-1">
              <el-tag v-for="p in (row.paths || [])" :key="p" size="small" effect="plain">{{ p }}</el-tag>
            </div>
          </template>
        </el-table-column>
        <el-table-column prop="category" :label="$t('auto.CompetitorManage.d0771a42')" width="100">
          <template #default="{ row }">
            <el-tag size="small" :type="categoryTag(row.category)">{{ categoryLabel(row.category) }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="priority" :label="$t('auto.CompetitorManage.ee8ecb9e')" width="80" sortable />
        <el-table-column prop="status" :label="$t('auto.CompetitorManage.3fea7ca7')" width="90">
          <template #default="{ row }">
            <el-tag size="small" :type="row.status === 'active' ? 'success' : 'info'">{{ row.status }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column :label="$t('auto.CompetitorManage.2b6bc0f2')" width="180">
          <template #default="{ row }">
            <el-button size="small" @click="openEdit(row)">{{ $t('auto.CompetitorManage.95b351c8') }}</el-button>
            <el-button size="small" type="danger" plain @click="onDelete(row)">{{ $t('auto.CompetitorManage.2f4aaddd') }}</el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    
    <el-dialog v-model="dialogVisible" ::title="$t('auto.CompetitorManage.e2045dc4')" width="520px" destroy-on-close>
      <el-form :model="form" label-width="90px" size="default">
        <el-form-item :label="$t('auto.CompetitorManage.d7ec2d3f')">
          <el-input v-model="form.name" :placeholder="$t('auto.CompetitorManage.68edbd0b')" />
        </el-form-item>
        <el-form-item :label="$t('auto.CompetitorManage.190980c7')">
          <el-input v-model="form.domain" :placeholder="$t('auto.CompetitorManage.0aa1c6a5')" :disabled="!!editing" />
        </el-form-item>
        <el-form-item :label="$t('auto.CompetitorManage.6147cfcd')">
          <div>
            <el-tag v-for="(p, i) in form.paths" :key="i" closable @close="form.paths.splice(i, 1)" size="default" style="margin-right:4px;margin-bottom:4px">
              {{ p }}
            </el-tag>
            <el-input v-if="newPath" v-model="newPath" size="small" style="width:120px" @keyup.enter="addPath" />
            <el-button v-else size="small" plain @click="newPath = '/'">{{ $t('auto.CompetitorManage.0470376f') }}</el-button>
          </div>
        </el-form-item>
        <el-form-item :label="$t('auto.CompetitorManage.d0771a42')">
          <el-select v-model="form.category">
            <el-option :label="$t('auto.CompetitorManage.710d58a4')" value="direct" />
            <el-option :label="$t('auto.CompetitorManage.80203f37')" value="global" />
            <el-option :label="$t('auto.CompetitorManage.3076094a')" value="indirect" />
          </el-select>
        </el-form-item>
        <el-form-item :label="$t('auto.CompetitorManage.ee8ecb9e')">
          <el-slider v-model="form.priority" :min="1" :max="10" show-input style="width:240px" />
        </el-form-item>
        <el-form-item :label="$t('auto.CompetitorManage.3fea7ca7')">
          <el-switch v-model="form.status" active-value="active" inactive-value="paused" active-text="active" inactive-text="paused" />
        </el-form-item>
        <el-form-item :label="$t('auto.CompetitorManage.2432b575')">
          <el-input v-model="form.notes" type="textarea" :rows="2" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">{{ $t('auto.CompetitorManage.625fb26b') }}</el-button>
        <el-button type="primary" :loading="saving" @click="onSave">{{ $t('auto.CompetitorManage.be5fbbe3') }}</el-button>
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

const categoryLabel = (c) => ({ direct: t('auto.CompetitorManage.710d58a4'), global: t('auto.CompetitorManage.80203f37'), indirect: t('auto.CompetitorManage.3076094a') }[c] || c || '—')
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
    ElMessage.warning(t('auto.CompetitorManage.59c49148'))
    return
  }
  saving.value = true
  try {
    if (editing.value) {
      await updateCompetitor(editing.value.id, form.value)
      ElMessage.success(t('auto.CompetitorManage.55aa6366'))
    } else {
      await createCompetitor(form.value)
      ElMessage.success(t('auto.CompetitorManage.04a691b3'))
    }
    dialogVisible.value = false
    await load()
  } catch (e) {
    ElMessage.error(e?.response?.data?.message || t('auto.CompetitorManage.6de920b4'))
  } finally {
    saving.value = false
  }
}

async function onDelete(row) {
  try {
    await ElMessageBox.confirm(`确认删除竞品 "${row.name}" 吗？`, t('auto.CompetitorManage.02d9819d'), { type: 'warning' })
    await deleteCompetitor(row.id)
    ElMessage.success(t('auto.CompetitorManage.5cc23262'))
    await load()
  } catch { /* 忽略：清理/存储/恢复类 best-effort 操作 */ }
}

onMounted(load)
</script>
