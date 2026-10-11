<template>
  <div class="app-container">
    <el-card shadow="never">
      <template #header>
        <div class="card-header">
          <span>{{ $t('auto.RBACMatrix.6.24a4a4') }}</span>
        </div>
      </template>

      <el-alert
        type="info"
        :closable="false"
        show-icon
        :title="$t('auto.RBACMatrix.1.e98201')"
        style="margin-bottom: 12px"
      />

      <template v-if="data">
        <el-alert
          type="warning"
          :closable="false"
          show-icon
          :title="$t('auto.RBACMatrix.2.a9db17')"
          style="margin-bottom: 12px"
        />
        <div style="margin-bottom: 12px">
          <el-button type="primary" size="small" :loading="savingMatrix" :disabled="!matrixDirty" @click="saveMatrix">{{ $t('auto.RBACMatrix.7.c5575b') }}</el-button>
          <el-button size="small" :disabled="!matrixDirty" @click="resetMatrix">{{ $t('auto.RBACMatrix.8.5ac201') }}</el-button>
        </div>
        <el-descriptions :column="4" border size="small" style="margin-bottom: 16px" :title="$t('auto.RBACMatrix.3.464f3d')">
          <el-descriptions-item v-for="r in data.roles" :key="r.id" :label="r.id">
            <el-tag size="small" :type="r.id === 'super_admin' ? 'danger' : r.id === 'team_admin' ? 'warning' : 'info'">
              {{ r.id }}
            </el-tag>
          </el-descriptions-item>
        </el-descriptions>

        <el-table :data="data.permissions" border size="small">
          <el-table-column prop="id" :label="$t('auto.RBACMatrix.4.2007ae')" min-width="200">
            <template #default="{ row }">
              <code class="mono">{{ row.id }}</code>
            </template>
          </el-table-column>
          <el-table-column
            v-for="r in data.roles"
            :key="r.id"
            :label="r.id"
            :width="roleColWidth"
            align="center"
          >
            <template #default="{ row }">
              <el-checkbox
                v-if="r.id !== 'super_admin'"
                :model-value="draft[r.id]?.includes(row.id) || false"
                @change="(v) => togglePerm(r.id, row.id, v)"
              />
              <span v-else>🔒</span>
            </template>
          </el-table-column>
        </el-table>

        <el-descriptions :column="4" border size="small" style="margin-top: 16px" :title="$t('auto.RBACMatrix.5.ff36d6')">
          <el-descriptions-item v-for="st in data.api_key_states" :key="st.id" :label="st.id">
            <code class="mono">{{ st.id }}</code>
          </el-descriptions-item>
        </el-descriptions>
      </template>
    </el-card>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { ElMessage } from 'element-plus'
import { http } from '@/utils/request'
import { rbacApi } from '@/api/teamAdmin'

const { t } = useI18n()

const data = ref(null)
// draft[roleCode] = [permIds]——仅覆盖可编辑角色（team_admin/member）；
// 值为 null 表示"无覆盖，跟生成矩阵走"，一旦勾选即产生覆盖。
const draft = reactive({})
const savingMatrix = ref(false)
const matrixDirty = computed(() => Object.keys(draft).length > 0)

const togglePerm = (roleCode, permId, checked) => {
  const base = draft[roleCode] !== undefined ? [...draft[roleCode]] : [...(data.value.permissions.filter((p) => p.roles.includes(roleCode)).map((p) => p.id))]
  const idx = base.indexOf(permId)
  if (checked && idx < 0) base.push(permId)
  if (!checked && idx >= 0) base.splice(idx, 1)
  draft[roleCode] = base
}

const saveMatrix = async () => {
  savingMatrix.value = true
  try {
    for (const [roleCode, perms] of Object.entries(draft)) {
      await http.put(`/api/manage/role-permissions/${roleCode}`, { permissions: perms })
    }
    ElMessage.success(t('auto.RBACMatrix.9.a7b722'))
    Object.keys(draft).forEach((k) => delete draft[k])
    await load()
  } catch (e) {
    ElMessage.error(e?.response?.data?.message || t('auto.RBACMatrix.10.6de920'))
  } finally {
    savingMatrix.value = false
  }
}

const resetMatrix = () => {
  Object.keys(draft).forEach((k) => delete draft[k])
}

const roleColWidth = computed(() => (data.value ? Math.floor(600 / (data.value.roles.length || 1)) : 120))

const load = async () => {
  const res = await rbacApi.matrix()
  data.value = res?.data ?? res ?? null
}

onMounted(load)
</script>

<style scoped>
.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}
.mono {
  font-family: monospace;
  font-size: 12px;
}
.muted {
  color: #ccc;
}
</style>
