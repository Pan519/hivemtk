<template>
  <div class="email-settings">
    <div class="toolbar">
      <el-button type="primary" @click="handleAdd">
        <el-icon-plus /> {{ $t('新增邮件代理') }}
      </el-button>
    </div>

    
    <el-table :data="tableData" border stripe style="width: 100%">
      <el-table-column prop="name" ::label="$t('auto.Smtp.77fbe7d8')" width="300" />
      <el-table-column prop="server" ::label="$t('auto.Smtp.b210fb15')" width="200" />
      <el-table-column prop="port" ::label="$t('auto.Smtp.2696c1d1')" width="100" />
      <el-table-column prop="limit" ::label="$t('auto.Smtp.6ae27167')" width="100" />
      <el-table-column prop="password" ::label="$t('auto.Smtp.b9ae4dc8')" />
      <el-table-column prop="username" ::label="$t('auto.Smtp.ef15518f')" />
      <el-table-column ::label="$t('auto.Smtp.4414e3c5')" width="200" align="center">
        <template #default="scope">
          <el-button size="small" @click="handleEdit(scope.row)">{{ $t('编辑') }}</el-button>
          <el-button size="small" type="danger" @click="handleDelete(scope.row.id)">{{ $t('删除') }}</el-button>
        </template>
      </el-table-column>
    </el-table>

    
    <el-dialog v-model="dialogVisible" :title="dialogTitle" width="600px">
      <el-form ref="formRef" :model="form" :rules="rules" label-width="120px">
        <el-form-item :label="$t('auto.Smtp.00d24bdf')" prop="name">
          <el-input v-model="form.name" :placeholder="$t('auto.Smtp.b6f633d7')" />
        </el-form-item>
        <el-form-item :label="$t('auto.Smtp.5d57821d')" prop="server">
          <el-input v-model="form.server" :placeholder="$t('auto.Smtp.c34c1bdc')" />
        </el-form-item>
        <el-form-item :label="$t('auto.Smtp.c76cfefe')" prop="port">
          <el-input v-model.number="form.port" :placeholder="$t('auto.Smtp.8d0fa2ee')" type="number" />
        </el-form-item>
        <el-form-item :label="$t('auto.Smtp.819767ad')" prop="username">
          <el-input v-model="form.username" :placeholder="$t('auto.Smtp.08b1fa13')" />
        </el-form-item>
        <el-form-item :label="$t('auto.Smtp.a8105204')" prop="password">
          <el-input v-model="form.password" :placeholder="$t('auto.Smtp.e39ffe99')" type="password" show-password autocomplete="new-password" />
        </el-form-item>
        <el-form-item :label="$t('auto.Smtp.d8f7d483')" prop="limit">
          <el-input v-model="form.limit" :placeholder="$t('auto.Smtp.742c9fd3')" type="number" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">{{ $t('auto.Smtp.625fb26b') }}</el-button>
        <el-button type="primary" @click="submitForm">{{ $t('auto.Smtp.38cf16f2') }}</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import i18n from '@/i18n'

import { ref, reactive, onMounted } from 'vue'
import { useI18n } from 'vue-i18n';
import { ElMessage, ElMessageBox } from 'element-plus';
import { Plus } from '@element-plus/icons-vue';

import { emailApi } from '@/api/email';

const { t } = useI18n()

const tableData = ref([]);

const fetchEmailSmtp = async () => {
  const res = await emailApi.getEmailSmtpList()
  tableData.value = res.list || [];
};

onMounted(() => {
  fetchEmailSmtp()
})


const dialogVisible = ref(false);
const dialogTitle = ref(t('auto.Smtp.247ae5e1'));
const formRef = ref(null);

const form = reactive({
  id: '',
  name: '',
  server: '',
  port: '',
  username: '',
  password: '',
  limit: 50
});

const rules = {
  name: [{ required: true, message: i18n.global.t(t('auto.Smtp.b6f633d7')), trigger: 'blur' }],
  server: [{ required: true, message: i18n.global.t(t('auto.Smtp.44056c2f')), trigger: 'blur' }],
  port: [
    { required: true, message: i18n.global.t(t('auto.Smtp.8d0fa2ee')), trigger: 'blur' },
    { type: 'number', message: i18n.global.t(t('auto.Smtp.7522a195')), trigger: 'blur' }
  ],
  username: [{ required: true, message: i18n.global.t(t('auto.Smtp.08b1fa13')), trigger: 'blur' }],
  password: [{ required: true, message: i18n.global.t(t('auto.Smtp.e39ffe99')), trigger: 'blur' }],
  limit: [{ required: true, message: i18n.global.t(t('auto.Smtp.742c9fd3')), trigger: 'blur' }],
};

const handleAdd = () => {
  dialogTitle.value = t('auto.Smtp.247ae5e1');
  Object.assign(form, { id: '', name: '', server: '', port: '', username: '', password: '', limit: 50 });
  dialogVisible.value = true;
};

const handleEdit = (row) => {
  dialogTitle.value = t('auto.Smtp.c1914748');
  Object.assign(form, { ...row });
  dialogVisible.value = true;
};

const handleDelete = (id) => {
  ElMessageBox.confirm(t('auto.Smtp.e27f8910'), t('auto.Smtp.900c70fa'), {
    confirmButtonText: t('auto.Smtp.38cf16f2'),
    cancelButtonText: t('auto.Smtp.625fb26b'),
    type: 'warning'
  }).then( async () => {
      try {
        const response = await emailApi.deleteEmailSmtp(id);
        ElMessage({
          type: 'success',
          message: i18n.global.t(t('auto.Smtp.0007d170')),
        });
        fetchEmailSmtp()
      } catch (error) {
        console.error(t('auto.Smtp.e4a3bb2d'), error)
        ElMessage.error(i18n.global.t(t('auto.Smtp.5c7c95b6')))
      }
  }).catch((e) => {
    if (e !== 'cancel' && e !== 'close')
      throw e;
  });
};

const submitForm = () => {
  form.limit = Number(form.limit)
  formRef.value.validate(async (valid) => {
    if (valid) {
          if(form.id){
            const res = await emailApi.updateEmailSmtp(form.id,form)
            ElMessage.success(i18n.global.t(t('auto.Smtp.55aa6366')));
            dialogVisible.value = false
            fetchEmailSmtp();
          }else{
            const res = await emailApi.addEmailSmtp(form)
            ElMessage.success(i18n.global.t(t('auto.Smtp.a5bfd70d')));
            dialogVisible.value = false
            fetchEmailSmtp();
          }
      dialogVisible.value = false;
    }
  });
};
</script>

<style scoped>
.toolbar {
  margin-bottom: 16px;
  display: flex;
  justify-content: flex-end;
}

.email-settings {
  padding: 10px;
}
</style>