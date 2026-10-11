<template>
  <div class="asset-market">
    <el-card class="filter-bar" shadow="never">
      <el-form inline>
        <el-form-item :label="$t('auto.Market.226b0912')">
          <el-select v-model="filter.asset_type" :placeholder="$t('auto.Market.a8b0c204')" clearable>
            <el-option :label="$t('auto.Market.5ce3a11b')" value="agent_persona" />
            <el-option :label="$t('auto.Market.655e227e')" value="sales_script" />
            <el-option :label="$t('auto.Market.b515e4c5')" value="ab_test_plan" />
            <el-option :label="$t('auto.Market.4a3b606a')" value="marketing_workflow" />
            <el-option :label="$t('auto.Market.112acd54')" value="industry_sop" />
          </el-select>
        </el-form-item>
        <el-form-item :label="$t('auto.Market.6adafe27')">
          <el-select v-model="filter.industry" :placeholder="$t('auto.Market.a8b0c204')" clearable>
            <el-option v-for="i in industries" :key="i" :label="i" :value="i" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="fetchList">{{ $t('auto.Market.bee912d7') }}</el-button>
          <el-button @click="$router.push('/asset-market/my-assets')">{{ $t('auto.Market.de43bfaf') }}</el-button>
        </el-form-item>
      </el-form>
    </el-card>

    <div class="asset-grid" v-loading="loading">
      <el-card
        v-for="a in list"
        :key="a.asset_id"
        class="asset-card"
        shadow="hover"
        @click="$router.push(`/asset-market/detail/${a.asset_id}`)"
      >
        <div class="cover-wrap">
          <img v-if="a.cover_url" :src="a.cover_url" class="cover" :alt="$t('auto.Market.5511d4ea')" />
          <div v-else class="cover-placeholder">{{ typeLabel(a.asset_type)[0] }}</div>
        </div>
        <h3 class="title">{{ a.name }}</h3>
        <div class="tags">
          <el-tag size="small" effect="plain">{{ a.industry }}</el-tag>
          <el-tag size="small" :type="typeColor(a.asset_type)">{{ typeLabel(a.asset_type) }}</el-tag>
        </div>
        <p class="desc">{{ a.description || '暂无描述' }}</p>
        <div class="footer">
          <div class="meta">
            <el-rate v-model="a.rating_avg" disabled size="small" />
            <span class="downloads">↓ {{ a.download_count || 0 }}</span>
          </div>
          <el-tag v-if="a.purchased" type="success" size="small">{{ $t('auto.Market.5af67104') }}</el-tag>
          <el-button v-else type="primary" size="small" @click.stop="handlePurchase(a)">
            {{ (a.price > 0) ? `${a.price} 积分购买` : '免费试用' }}
          </el-button>
        </div>
      </el-card>
    </div>

    <el-empty v-if="!loading && list.length === 0" :description="$t('auto.Market.1b91cf96')" />

    <el-pagination
      v-model:current-page="filter.page"
      v-model:page-size="filter.size"
      :total="total"
      layout="total, prev, pager, next"
      @current-change="fetchList"
    />
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { listAssets, purchaseAsset } from '@/api/assetMarket'
import { INDUSTRY_VALUES } from '@/constants/industry'
import { ElMessage, ElMessageBox } from 'element-plus'

const industries = INDUSTRY_VALUES
const { t } = useI18n()

const filter = ref({ asset_type: '', industry: '', page: 1, size: 20 })
const list = ref([])
const total = ref(0)
const loading = ref(false)

const typeLabel = (t) =>
  ({
    agent_persona: t('auto.Market.5ce3a11b'),
    sales_script: t('auto.Market.655e227e'),
    ab_test_plan: 'AB 测试',
    marketing_workflow: t('auto.Market.228526dc'),
    industry_sop: t('auto.Market.112acd54')
  }[t] || t)

const typeColor = (t) =>
  ({
    agent_persona: 'primary',
    sales_script: 'success',
    ab_test_plan: 'warning',
    marketing_workflow: 'info',
    industry_sop: 'danger'
  }[t] || '')

const fetchList = async () => {
  loading.value = true
  try {
    const resp = await listAssets(filter.value)
    const data = resp?.data || resp || {}
    list.value = data.list || data || []
    total.value = data.total || 0
  } finally {
    loading.value = false
  }
}

const handlePurchase = async (asset) => {
  const paid = (asset.price > 0)
  const confirmMsg = paid
    ? `确认以 ${asset.price} 积分购买资产「${asset.name}」？购买后将自动同步到本地（积分余额不足将被拒绝）。`
    : `确认「免费试用」资产「${asset.name}」？试用后将自动同步到本地。`
  try {
    await ElMessageBox.confirm(confirmMsg, paid ? t('auto.Market.f1195e12') : t('auto.Market.248f81e2'), { type: 'warning' })
  } catch {
    return
  }
  await purchaseAsset({ asset_id: asset.asset_id })
  ElMessage.success(paid ? t('auto.Market.b5442a8a') : t('auto.Market.e8d3b5ac'))
  fetchList()
}

onMounted(fetchList)
</script>

<style scoped lang="scss">
.asset-market {
  padding: 16px;
  .filter-bar {
    margin-bottom: 16px;
  }
  .asset-grid {
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(280px, 1fr));
    gap: 16px;
  }
  .asset-card {
    cursor: pointer;
    transition: transform 0.2s;
    &:hover {
      transform: translateY(-4px);
    }
    .cover-wrap {
      height: 140px;
      overflow: hidden;
      border-radius: 4px;
      margin-bottom: 12px;
      .cover {
        width: 100%;
        height: 100%;
        object-fit: cover;
      }
      .cover-placeholder {
        width: 100%;
        height: 100%;
        display: flex;
        align-items: center;
        justify-content: center;
        background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
        color: white;
        font-size: 48px;
        font-weight: bold;
      }
    }
    .title {
      font-size: 16px;
      margin: 0 0 8px;
    }
    .tags {
      margin-bottom: 8px;
      .el-tag {
        margin-right: 4px;
      }
    }
    .desc {
      color: #666;
      font-size: 13px;
      min-height: 40px;
      margin: 0 0 12px;
    }
    .footer {
      display: flex;
      justify-content: space-between;
      align-items: center;
    }
    .meta {
      display: flex;
      align-items: center;
      gap: 8px;
      color: #999;
      font-size: 12px;
    }
  }
  .el-pagination {
    margin-top: 16px;
    justify-content: center;
  }
}
</style>
