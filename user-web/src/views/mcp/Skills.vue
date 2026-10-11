<template>
  <div class="app-container">
    <el-card shadow="never">
      <template #header>
        <div class="card-header">
          <span>{{ $t('auto.Skills.4c8ae944') }}</span>
          <el-button type="primary" size="small" @click="openCreate">{{ $t('auto.Skills.ee41104e') }}</el-button>
        </div>
      </template>

      <el-alert
        type="info"
        :closable="false"
        show-icon
        :title="$t('auto.Skills.4ab32071')"
        style="margin-bottom: 12px"
      />

      <el-table v-loading="loading" :data="rows" border size="small">
        <el-table-column prop="name" :label="$t('auto.Skills.d7ec2d3f')" width="160" />
        <el-table-column prop="description" :label="$t('auto.Skills.f411d0f1')" min-width="200" show-overflow-tooltip />
        <el-table-column :label="$t('auto.Skills.b98c84eb')" width="140">
          <template #default="{ row }">
            <span v-if="row.credential_id">{{ '凭证 #' + row.credential_id }}</span>
            <el-tag v-else type="info" size="small">{{ $t('auto.Skills.906ad18b') }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column :label="$t('auto.Skills.d1f4a2ca')" min-width="180">
          <template #default="{ row }">
            <el-tag v-for="t in (row.allowed_tools || '').split(',').filter(Boolean)" :key="t" size="small" style="margin-right: 4px">
              {{ t }}
            </el-tag>
            <span v-if="!row.allowed_tools" class="muted">{{ $t('auto.Skills.c9fc9997') }}</span>
          </template>
        </el-table-column>
        <el-table-column :label="$t('auto.Skills.3fea7ca7')" width="80">
          <template #default="{ row }">
            <el-tag :type="row.enabled ? 'success' : 'info'" size="small">{{ row.enabled ? '启用' : '停用' }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column :label="$t('auto.Skills.2b6bc0f2')" width="200" fixed="right">
          <template #default="{ row }">
            <el-button size="small" @click="openEdit(row)">{{ $t('auto.Skills.95b351c8') }}</el-button>
            <el-button size="small" @click="toggle(row)">{{ row.enabled ? '停用' : '启用' }}</el-button>
            <el-button size="small" type="danger" @click="remove(row)">{{ $t('auto.Skills.2f4aaddd') }}</el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <!-- 登记 / 编辑弹窗 -->
    <el-dialog v-model="editorVisible" ::title="$t('auto.Skills.a63a1c87')" width="560px">
      <el-form label-width="90px">
        <el-form-item :label="$t('auto.Skills.d7ec2d3f')" required>
          <el-input v-model="form.name" :placeholder="$t('auto.Skills.4f3a5228')" maxlength="128" />
        </el-form-item>
        <el-form-item :label="$t('auto.Skills.f411d0f1')">
          <el-input v-model="form.description" type="textarea" :rows="2" :placeholder="$t('auto.Skills.a04f8972')" />
        </el-form-item>
        <el-form-item :label="$t('auto.Skills.b98c84eb')">
          <el-select v-model="form.credential_id" clearable :placeholder="$t('auto.Skills.2239705b')" style="width: 100%">
            <el-option v-for="c in credentials" :key="c.id" :label="c.client_id" :value="c.id" />
          </el-select>
        </el-form-item>
        <el-form-item :label="$t('auto.Skills.d1f4a2ca')">
          <el-input v-model="form.allowed_tools" :placeholder="$t('auto.Skills.da5711d8')" />
        </el-form-item>
        <el-form-item :label="$t('auto.Skills.270c7dfc')">
          <el-input v-model="form.doc_url" :placeholder="$t('auto.Skills.adc81568')" />
        </el-form-item>
        <el-form-item :label="$t('auto.Skills.c86b6062')">
          <el-input v-model="form.guardrails" type="textarea" :rows="2" :placeholder="$t('auto.Skills.72a47c53')" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="editorVisible = false">{{ $t('auto.Skills.625fb26b') }}</el-button>
        <el-button type="primary" :loading="saving" @click="save">{{ $t('auto.Skills.be5fbbe3') }}</el-button>
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
    ElMessage.warning(t('auto.Skills.856d0df5'))
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
      ElMessage.success(t('auto.Skills.1b92092d'))
    } else {
      await mcpSkillApi.create(payload)
      ElMessage.success(t('auto.Skills.689286b4'))
    }
    editorVisible.value = false
    await load()
  } catch (e) {
    ElMessage.error(e?.response?.data?.message || t('auto.Skills.6de920b4'))
  } finally {
    saving.value = false
  }
}

const toggle = async (row) => {
  await mcpSkillApi.update(row.id, { enabled: !row.enabled })
  ElMessage.success(row.enabled ? t('auto.Skills.69b0f684') : t('auto.Skills.53ace430'))
  await load()
}

const remove = async (row) => {
  await ElMessageBox.confirm(`确认删除 Skill ${row.name}？`, t('auto.Skills.50eaf94d'), { type: 'warning' })
  await mcpSkillApi.remove(row.id)
  ElMessage.success(t('auto.Skills.5cc23262'))
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