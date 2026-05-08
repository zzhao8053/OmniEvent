<template>
  <v-container class="fill-height" fluid>
    <v-row align="center" justify="center">
      <v-col cols="12" sm="8" md="6">
        <v-card class="elevation-12">
          <v-card-title class="text-center">
            <h2>个人资料</h2>
          </v-card-title>
          <v-card-text>
            <v-form ref="formRef" @submit.prevent="handleUpdate">
              <v-text-field
                v-model="profile.username"
                label="用户名"
                variant="outlined"
                prepend-inner-icon="mdi-account"
                disabled
              />
              <v-text-field
                v-model="profile.nickname"
                label="昵称"
                :rules="[rules.required]"
                variant="outlined"
                prepend-inner-icon="mdi-card-account-details"
              />
              <v-text-field
                v-model="profile.email"
                label="邮箱"
                :rules="[rules.required, rules.email]"
                variant="outlined"
                prepend-inner-icon="mdi-email"
              />
              <v-alert v-if="errorMessage" type="error" class="mb-4">
                {{ errorMessage }}
              </v-alert>
              <v-alert v-if="successMessage" type="success" class="mb-4">
                {{ successMessage }}
              </v-alert>
              <v-btn
                type="submit"
                color="primary"
                block
                :loading="loading"
                size="large"
              >
                保存修改
              </v-btn>
            </v-form>
          </v-card-text>
          <v-card-actions>
            <v-spacer />
            <v-btn color="error" variant="text" @click="handleLogout">
              退出登录
            </v-btn>
          </v-card-actions>
        </v-card>
      </v-col>
    </v-row>
  </v-container>
</template>

<script setup lang="ts">
import { ref, onMounted, reactive } from 'vue'
import { useRouter } from 'vue-router'
import { useUserStore } from '@/stores/user'

const router = useRouter()
const userStore = useUserStore()

const profile = reactive({
  username: '',
  nickname: '',
  email: ''
})

const loading = ref(false)
const errorMessage = ref('')
const successMessage = ref('')
const formRef = ref<any>(null)

const rules = {
  required: (v: string) => !!v || '此字段为必填项',
  email: (v: string) => /.+@.+\..+/.test(v) || '请输入有效的邮箱地址'
}

onMounted(async () => {
  try {
    const data = await userStore.fetchProfile()
    if (data) {
      profile.username = data.username
      profile.nickname = data.nickname
      profile.email = data.email
    }
  } catch (error) {
    console.error('Failed to fetch profile:', error)
  }
})

async function handleUpdate() {
  const { valid } = await formRef.value.validate()
  if (!valid) return

  loading.value = true
  errorMessage.value = ''
  successMessage.value = ''

  try {
    await userStore.updateProfile(profile.nickname, profile.email)
    successMessage.value = '资料更新成功！'
  } catch (error: any) {
    errorMessage.value = error.response?.data?.error || '更新失败，请稍后重试'
  } finally {
    loading.value = false
  }
}

function handleLogout() {
  userStore.logout()
  router.push('/login')
}
</script>
