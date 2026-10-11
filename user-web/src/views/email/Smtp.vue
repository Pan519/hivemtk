<template>
  <div class="email-settings">
    <div class="toolbar">
      <el-button type="primary" @click="handleAdd">
        <el-icon-plus /> {{ $t('新增邮件代理') }}
      </el-button>
    </div>

    
    <el-table :data="tableData" border stripe style="width: 100%">
      <el-table-column prop="name" ::label="$t('auto.Smtp.1.77fbe7')" width="300" />
      <el-table-column prop="server" ::label="$t('auto.Smtp.2.b210fb')" width="200" />
      <el-table-column prop="port" ::label="$t('auto.Smtp.3.2696c1')" width="100" />
      <el-table-column prop="limit" ::label="$t('auto.Smtp.4.6ae271')" width="100" />
      <el-table-column prop="password" ::label="$t('auto.Smtp.5.b9ae4d')" />
      <el-table-column prop="username" ::label="$t('auto.Smtp.6.ef1551')" />
      <el-table-column ::label="$t('auto.Smtp.7.4414e3')" width="200" align="center">
        <template #default="scope">
          <el-button size="small" @click="handleEdit(scope.row)">{{ $t('编辑') }}</el-button>
          <el-button size="small" type="danger" @click="handleDelete(scope.row.id)">{{ $t('删除') }}</el-button>
        </template>
      </el-table-column>
    </el-table>

    
    <el-dialog v-model="dialogVisible" :title="dialogTitle" width="600px">
      <el-form ref="formRef" :model="form" :rules="rules" label-width="120px">
        <el-form-item :label="$t('auto.Smtp.8.00d24b')" prop="name">
          <el-input v-model="form.name" :placeholder="$t('auto.Smtp.9.b6f633')" />
        </el-form-item>
        <el-form-item :label="$t('auto.Smtp.10.5d5782')" prop="server">
          <el-input v-model="form.server" :placeholder="$t('auto.Smtp.11.c34c1b')" />
        </el-form-item>
        <el-form-item :label="$t('auto.Smtp.12.c76cfe')" prop="port">
          <el-input v-model.number="form.port" :placeholder="$t('auto.Smtp.13.8d0fa2')" type="number" />
        </el-form-item>
        <el-form-item :label="$t('auto.Smtp.14.819767')" prop="username">
          <el-input v-model="form.username" :placeholder="$t('auto.Smtp.15.08b1fa')" />
        </el-form-item>
        <el-form-item :label="$t('auto.Smtp.16.a81052')" prop="password">
          <el-input v-model="form.password" :placeholder="$t('auto.Smtp.17.e39ffe')" type="password" show-password autocomplete="new-password" />
        </el-form-item>
        <el-form-item :label="$t('auto.Smtp.18.d8f7d4')" prop="limit">
          <el-input v-model="form.limit" :placeholder="$t('auto.Smtp.19.742c9f')" type="number" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">{{ $t('auto.Smtp.20.625fb2') }}</el-button>
        <el-button type="primary" @click="submitForm">{{ $t('auto.Smtp.21.38cf16') }}</el-button>
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
const dialogTitle = ref(t('auto.Smtp.22.247ae5'));
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
  name: [{ required: true, message: i18n.global.t(t('auto.Smtp.9.b6f633')), trigger: 'blur' }],
  server: [{ required: true, message: i18n.global.t(t('auto.Smtp.23.44056c')), trigger: 'blur' }],
  port: [
    { required: true, message: i18n.global.t(t('auto.Smtp.13.8d0fa2')), trigger: 'blur' },
    { type: 'number', message: i18n.global.t(t('auto.Smtp.24.7522a1')), trigger: 'blur' }
  ],
  username: [{ required: true, message: i18n.global.t(t('auto.Smtp.15.08b1fa')), trigger: 'blur' }],
  password: [{ required: true, message: i18n.global.t(t('auto.Smtp.17.e39ffe')), trigger: 'blur' }],
  limit: [{ required: true, message: i18n.global.t(t('auto.Smtp.19.742c9f')), trigger: 'blur' }],
};

const handleAdd = () => {
  dialogTitle.value = t('auto.Smtp.22.247ae5');
  Object.assign(form, { id: '', name: '', server: '', port: '', username: '', password: '', limit: 50 });
  dialogVisible.value = true;
};

const handleEdit = (row) => {
  dialogTitle.value = t('auto.Smtp.25.c19147');
  Object.assign(form, { ...row });
  dialogVisible.value = true;
};

const handleDelete = (id) => {
  ElMessageBox.confirm(t('auto.Smtp.26.e27f89'), t('auto.Smtp.27.900c70'), {
    confirmButtonText: t('auto.Smtp.21.38cf16'),
    cancelButtonText: t('auto.Smtp.20.625fb2'),
    type: 'warning'
  }).then( async () => {
      try {
        const response = await emailApi.deleteEmailSmtp(id);
        ElMessage({
          type: 'success',
          message: i18n.global.t(t('auto.Smtp.28.0007d1')),
        });
        fetchEmailSmtp()
      } catch (error) {
        console.error(t('auto.Smtp.29.e4a3bb'), error)
        ElMessage.error(i18n.global.t(t('auto.Smtp.30.5c7c95')))
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
            ElMessage.success(i18n.global.t(t('auto.Smtp.31.55aa63')));
            dialogVisible.value = false
            fetchEmailSmtp();
          }else{
            const res = await emailApi.addEmailSmtp(form)
            ElMessage.success(i18n.global.t(t('auto.Smtp.32.a5bfd7')));
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