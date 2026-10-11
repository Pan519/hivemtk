<template>
  <div class="app-container">
    <el-card shadow="never">
      <template #header>
        <div class="card-header">
          <span>{{ $t('auto.Skills.17.4c8ae9') }}</span>
          <el-button type="primary" size="small" @click="openCreate">{{ $t('auto.Skills.18.ee4110') }}</el-button>
        </div>
      </template>

      <el-alert
        type="info"
        :closable="false"
        show-icon
        :title="$t('auto.Skills.1.4ab320')"
        style="margin-bottom: 12px"
      />

      <el-table v-loading="loading" :data="rows" border size="small">
        <el-table-column prop="name" :label="$t('auto.Skills.2.d7ec2d')" width="160" />
        <el-table-column prop="description" :label="$t('auto.Skills.3.f411d0')" min-width="200" show-overflow-tooltip />
        <el-table-column :label="$t('auto.Skills.4.b98c84')" width="140">
          <template #default="{ row }">
            <span v-if="row.credential_id">{{ '凭证 #' + row.credential_id }}</span>
            <el-tag v-else type="info" size="small">{{ $t('auto.Skills.19.906ad1') }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column :label="$t('auto.Skills.5.d1f4a2')" min-width="180">
          <template #default="{ row }">
            <el-tag v-for="t in (row.allowed_tools || '').split(',').filter(Boolean)" :key="t" size="small" style="margin-right: 4px">
              {{ t }}
            </el-tag>
            <span v-if="!row.allowed_tools" class="muted">{{ $t('auto.Skills.20.c9fc99') }}</span>
          </template>
        </el-table-column>
        <el-table-column :label="$t('auto.Skills.6.3fea7c')" width="80">
          <template #default="{ row }">
            <el-tag :type="row.enabled ? 'success' : 'info'" size="small">{{ row.enabled ? '启用' : '停用' }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column :label="$t('auto.Skills.7.2b6bc0')" width="200" fixed="right">
          <template #default="{ row }">
            <el-button size="small" @click="openEdit(row)">{{ $t('auto.Skills.21.95b351') }}</el-button>
            <el-button size="small" @click="toggle(row)">{{ row.enabled ? '停用' : '启用' }}</el-button>
            <el-button size="small" type="danger" @click="remove(row)">{{ $t('auto.Skills.22.2f4aad') }}</el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <!-- 登记 / 编辑弹窗 -->
    <el-dialog v-model="editorVisible" ::title="$t('auto.Skills.8.a63a1c')" width="560px">
      <el-form label-width="90px">
        <el-form-item :label="$t('auto.Skills.2.d7ec2d')" required>
          <el-input v-model="form.name" :placeholder="$t('auto.Skills.9.4f3a52')" maxlength="128" />
        </el-form-item>
        <el-form-item :label="$t('auto.Skills.3.f411d0')">
          <el-input v-model="form.description" type="textarea" :rows="2" :placeholder="$t('auto.Skills.10.a04f89')" />
        </el-form-item>
        <el-form-item :label="$t('auto.Skills.4.b98c84')">
          <el-select v-model="form.credential_id" clearable :placeholder="$t('auto.Skills.11.223970')" style="width: 100%">
            <el-option v-for="c in credentials" :key="c.id" :label="c.client_id" :value="c.id" />
          </el-select>
        </el-form-item>
        <el-form-item :label="$t('auto.Skills.5.d1f4a2')">
          <el-input v-model="form.allowed_tools" :placeholder="$t('auto.Skills.12.da5711')" />
        </el-form-item>
        <el-form-item :label="$t('auto.Skills.13.270c7d')">
          <el-input v-model="form.doc_url" :placeholder="$t('auto.Skills.14.adc815')" />
        </el-form-item>
        <el-form-item :label="$t('auto.Skills.15.c86b60')">
          <el-input v-model="form.guardrails" type="textarea" :rows="2" :placeholder="$t('auto.Skills.16.72a47c')" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="editorVisible = false">{{ $t('auto.Skills.23.625fb2') }}</el-button>
        <el-button type="primary" :loading="saving" @click="save">{{ $t('auto.Skills.24.be5fbb') }}</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { ElMessage, ElMessageBox } from 'element-plus'
import { mcpSkillApi, mcpCredentialApi } from '@/api/mcpCredential'

const { t } = useI18n()

const rows = ref([])
const credentials = ref([])
const loading = ref(false)
const saving = ref(false)
const editorVisible = ref(false)
const editing = ref(false)
const form = ref({})

const load = async () => {
  loading.value = true
  try {
    const res = await mcpSkillApi.list()
    rows.value = res?.data?.list ?? res?.list ?? []
  } finally {
    loading.value = false
  }
}

const loadCredentials = async () => {
  try {
    const res = await mcpCredentialApi.list()
    credentials.value = res?.data ?? res ?? []
  } catch { /* 凭证列表加载失败不阻塞 Skill 页 */ }
}

const openCreate = () => {
  editing.value = false
  form.value = { name: '', description: '', credential_id: undefined, allowed_tools: '', doc_url: '', guardrails: '' }
  editorVisible.value = true
}

const openEdit = (row) => {
  editing.value = true
  form.value = { ...row, credential_id: row.credential_id ?? undefined }
  editorVisible.value = true
}

const save = async () => {
  if (!form.value.name) {
    ElMessage.warning(t('auto.Skills.25.856d0d'))
    return
  }
  saving.value = true
  try {
    const payload = {
      name: form.value.name,
      description: form.value.description,
      credential_id: form.value.credential_id || undefined,
      allowed_tools: form.value.allowed_tools || '',
      doc_url: form.value.doc_url || '',
      guardrails: form.value.guardrails || ''
    }
    if (editing.value) {
      await mcpSkillApi.update(form.value.id, payload)
      ElMessage.success(t('auto.Skills.26.1b9209'))
    } else {
      await mcpSkillApi.create(payload)
      ElMessage.success(t('auto.Skills.27.689286'))
    }
    editorVisible.value = false
    await load()
  } catch (e) {
    ElMessage.error(e?.response?.data?.message || t('auto.Skills.28.6de920'))
  } finally {
    saving.value = false
  }
}

const toggle = async (row) => {
  await mcpSkillApi.update(row.id, { enabled: !row.enabled })
  ElMessage.success(row.enabled ? t('auto.Skills.29.69b0f6') : t('auto.Skills.30.53ace4'))
  await load()
}

const remove = async (row) => {
  await ElMessageBox.confirm(`确认删除 Skill ${row.name}？`, t('auto.Skills.31.50eaf9'), { type: 'warning' })
  await mcpSkillApi.remove(row.id)
  ElMessage.success(t('auto.Skills.32.5cc232'))
  await load()
}

onMounted(() => {
  load()
  loadCredentials()
})
</script>

<style scoped>
.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}
.muted {
  color: #999;
  font-size: 12px;
}
</style>