<template>
  <el-dialog
    :model-value="modelValue"
    :title="$t('auto.BlacklistDialog.5adc9f43')"
    width="680px"
    :close-on-click-modal="false"
    @update:model-value="$emit('update:modelValue', $event)"
  >
    <div class="blacklist-toolbar">
      <span class="count">共 {{ blacklistItems.length }} 条记录</span>
      <el-button size="small" :loading="blacklistLoading" @click="$emit('reload')">{{ $t('auto.BlacklistDialog.694fc5ef') }}</el-button>
    </div>
    <el-table
      v-loading="blacklistLoading"
      :data="blacklistItems"
      size="small"
      max-height="420"
      :empty-text="$t('auto.BlacklistDialog.dc5b46a1')"
    >
      <el-table-column :label="$t('auto.BlacklistDialog.79b519d4')" min-width="130">
        <template #default="{ row }">{{ row.user_id ?? row.UserID ?? row.userId ?? '-' }}</template>
      </el-table-column>
      <el-table-column :label="$t('auto.BlacklistDialog.4981287f')" width="100">
        <template #default="{ row }">{{ getChannelLabel(row.platform ?? row.Platform) }}</template>
      </el-table-column>
      <el-table-column :label="$t('auto.BlacklistDialog.67dc0fb6')" min-width="160" show-overflow-tooltip>
        <template #default="{ row }">{{ row.reason ?? row.Reason ?? '-' }}</template>
      </el-table-column>
      <el-table-column :label="$t('auto.BlacklistDialog.26ca20b1')" width="90">
        <template #default="{ row }">
          <el-tag size="small" :type="(row.source ?? row.Source) === 'auto' ? 'warning' : 'info'">
            {{ row.source ?? row.Source ?? 'manual' }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column :label="$t('auto.BlacklistDialog.cf252d0f')" width="140">
        <template #default="{ row }">{{ formatTime(row.created_at ?? row.CreatedAt) }}</template>
      </el-table-column>
      <el-table-column :label="$t('auto.BlacklistDialog.2b6bc0f2')" width="90" fixed="right">
        <template #default="{ row }">
          <el-button size="small" type="danger" link @click="$emit('unblacklist', row)">{{ $t('auto.BlacklistDialog.2d92531b') }}</el-button>
        </template>
      </el-table-column>
    </el-table>
    <template #footer>
      <el-button @click="$emit('update:modelValue', false)">{{ $t('auto.BlacklistDialog.b15d9127') }}</el-button>
    </template>
  </el-dialog>
</template>

<script setup>
// 黑名单管理弹窗(由 List.vue 原样迁出,渲染与交互保持一致)
import { getChannelLabel } from '@/constants/channel'
import { formatTime } from '../composables/useSessionFilters'

defineProps({
  modelValue: { type: Boolean, default: false },
  blacklistItems: { type: Array, default: () => [] },
  blacklistLoading: { type: Boolean, default: false }
})

defineEmits(['update:modelValue', 'reload', 'unblacklist'])
</script>

<style scoped>
/* C2：黑名单管理弹窗工具栏 */
.blacklist-toolbar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 12px;
}
.blacklist-toolbar .count {
  font-size: 13px;
  color: #606266;
}
</style>
