<template>
  <div class="session-tag-page">
    <el-card class="header-card">
      <div>
        <h2>{{ $t('会话标签') }}</h2>
        <p class="subtitle">{{ $t('维护客服会话分类标签，用于会话打标、复盘与统计') }}</p>
      </div>
      <el-button type="primary" @click="openCreate">
        <el-icon><Plus /></el-icon>
        {{ $t('新增标签') }}
      </el-button>
    </el-card>

    <el-card>
      <template #header>
        <div class="card-header">
          <el-space>
            <el-select v-model="filterGroup" ::placeholder="$t('auto.SessionTag.2809a00e')" clearable style="width: 180px">
              <el-option v-for="g in groups" :key="g" :label="g" :value="g" />
            </el-select>
            <el-input v-model="search" ::placeholder="$t('auto.SessionTag.24731aed')" clearable style="width: 220px" />
          </el-space>
        </div>
      </template>

      <el-table :data="filtered" v-loading="loading" stripe>
        <el-table-column prop="name" :label="$t('auto.SessionTag.14d34236')" min-width="140">
          <template #default="{ row }">
            <el-tag :color="row.color" effect="dark" style="color: #fff">
              {{ row.name }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="code" :label="$t('auto.SessionTag.45235f6c')" min-width="140" />
        <el-table-column prop="group" :label="$t('auto.SessionTag.829abe5a')" width="140" />
        <el-table-column prop="description" :label="$t('auto.SessionTag.f411d0f1')" min-width="200" show-overflow-tooltip />
        <el-table-column prop="sort_order" :label="$t('auto.SessionTag.c360e994')" width="80" align="center" />
        <el-table-column :label="$t('auto.SessionTag.2b6bc0f2')" width="180" fixed="right">
          <template #default="{ row }">
            <el-button link type="primary" @click="openEdit(row)">{{ $t('auto.SessionTag.95b351c8') }}</el-button>
            <el-button link type="danger" @click="handleDelete(row)">{{ $t('auto.SessionTag.2f4aaddd') }}</el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <el-dialog v-model="dialogVisible" ::title="$t('auto.SessionTag.334671d0')" width="480px">
      <el-form :model="form" label-width="100px">
        <el-form-item :label="$t('auto.SessionTag.0718bfa1')" required>
          <el-input v-model="form.name" />
        </el-form-item>
        <el-form-item :label="$t('auto.SessionTag.45235f6c')" required>
          <el-input v-model="form.code" :disabled="isEdit" :placeholder="$t('auto.SessionTag.ac684da7')" />
        </el-form-item>
        <el-form-item :label="$t('auto.SessionTag.829abe5a')">
          <el-input v-model="form.group" :placeholder="$t('auto.SessionTag.525acdad')" />
        </el-form-item>
        <el-form-item :label="$t('auto.SessionTag.6b36c6f7')">
          <el-color-picker v-model="form.color" />
        </el-form-item>
        <el-form-item :label="$t('auto.SessionTag.f411d0f1')">
          <el-input v-model="form.description" type="textarea" :rows="2" />
        </el-form-item>
        <el-form-item :label="$t('auto.SessionTag.c360e994')">
          <el-input-number v-model="form.sort" :min="0" :max="999" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">{{ $t('auto.SessionTag.625fb26b') }}</el-button>
        <el-button type="primary" @click="submit">{{ $t('auto.SessionTag.be5fbbe3') }}</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import i18n from '@/i18n'

import { ref, reactive, computed, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Plus } from '@element-plus/icons-vue'
import {
  getSessionTags,
  createSessionTag,
  updateSessionTag,
  deleteSessionTag
} from '@/api/customerService.js'

const { t } = useI18n()

const loading = ref(false)
const search = ref('')
const filterGroup = ref('')
const tags = ref([])

const dialogVisible = ref(false)
const isEdit = ref(false)
const form = reactive({
  id: null,
  name: '',
  code: '',
  group: '',
  color: '#4F46E5',
  description: '',
  sort: 0
})

const groups = computed(() => [...new Set(tags.value.map((t) => t.group).filter(Boolean))])

const filtered = computed(() => {
  let list = tags.value
  if (filterGroup.value) list = list.filter((t) => t.group === filterGroup.value)
  if (search.value) {
    const kw = search.value.toLowerCase()
    list = list.filter((t) => t.name?.toLowerCase().includes(kw) || t.code?.toLowerCase().includes(kw))
  }
  return list
})

const load = async () => {
  loading.value = true
  try {
    const res = await getSessionTags()
    const data = res || []
    tags.value = Array.isArray(data) ? data : data.list || []
  } catch (e) {
    ElMessage.error(i18n.global.t(t('auto.SessionTag.d895f912')))
    tags.value = []
  } finally {
    loading.value = false
  }
}

const openCreate = () => {
  isEdit.value = false
  Object.assign(form, {
    id: null, name: '', code: '', group: '', color: '#4F46E5', description: '', sort: 0
  })
  dialogVisible.value = true
}

const openEdit = (row) => {
  isEdit.value = true
  Object.assign(form, { ...row })
  dialogVisible.value = true
}

const submit = async () => {
  if (!form.name || !form.code) {
    ElMessage.warning(i18n.global.t(t('auto.SessionTag.c327e3e0')))
    return
  }
  const payload = {
    name: form.name,
    code: form.code,
    group: form.group,
    color: form.color,
    description: form.description,
    sort_order: form.sort
  }
  try {
    if (isEdit.value) {
      await updateSessionTag(form.id, payload)
      ElMessage.success(i18n.global.t(t('auto.SessionTag.1b92092d')))
    } else {
      await createSessionTag(payload)
      ElMessage.success(i18n.global.t(t('auto.SessionTag.2d4c5dcc')))
    }
    dialogVisible.value = false
    await load()
  } catch (e) {
    ElMessage.error(e?.message || t('auto.SessionTag.5fa802be'))
  }
}

const handleDelete = async (row) => {
  try {
    await ElMessageBox.confirm(`确定删除标签 "${row.name}"？`, t('auto.SessionTag.e83a256e'), { type: 'warning' })
    await deleteSessionTag(row.id)
    ElMessage.success(i18n.global.t(t('auto.SessionTag.0007d170')))
    await load()
  } catch (e) {
    if (e !== 'cancel' && e !== 'close') ElMessage.error(i18n.global.t(t('auto.SessionTag.acf0664a')))
  }
}

onMounted(() => load())
</script>

<style scoped lang="scss">
.session-tag-page { padding: 20px; }
.header-card {
  margin-bottom: 20px;
  :deep(.el-card__body) { display: flex; justify-content: space-between; align-items: center; }
  h2 { margin: 0 0 8px 0; }
  .subtitle { color: #909399; margin: 0; }
}
.card-header { display: flex; justify-content: space-between; align-items: center; }
</style>
