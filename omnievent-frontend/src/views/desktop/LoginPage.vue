<template>
  <v-container class="fill-height" fluid>
    <v-row align="center" justify="center">
      <v-col cols="12" sm="8" md="4">
        <v-card class="elevation-12">
          <v-card-title class="text-center">
            <h2>登录 OmniEvent</h2>
          </v-card-title>
          <v-card-text>
            <v-form ref="formRef" @submit.prevent="handleLogin">
              <v-text-field
                v-model="username"
                label="用户名或邮箱"
                :rules="[rules.required]"
                variant="outlined"
                prepend-inner-icon="mdi-account"
              />
              <v-text-field
                v-model="password"
                label="密码"
                :type="showPassword ? 'text' : 'password'"
                :rules="[rules.required]"
                variant="outlined"
                prepend-inner-icon="mdi-lock"
                :append-inner-icon="showPassword ? 'mdi-eye-off' : 'mdi-eye'"
                @click:append-inner="showPassword = !showPassword"
              />
              <v-alert v-if="errorMessage" type="error" class="mb-4">
                {{ errorMessage }}
              </v-alert>
              <v-btn
                type="submit"
                color="primary"
                block
                :loading="loading"
                size="large"
              >
                登录
              </v-btn>
            </v-form>
          </v-card-text>
          <v-card-actions>
            <v-spacer />
            <v-btn variant="text" @click="$router.push('/signup')">
              还没有账号？去注册
            </v-btn>
          </v-card-actions>
        </v-card>
      </v-col>
    </v-row>
  </v-container>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useUserStore } from '@/stores/user'

const router = useRouter()
const userStore = useUserStore()

const username = ref('')
const password = ref('')
const showPassword = ref(false)
const loading = ref(false)
const errorMessage = ref('')
const formRef = ref<any>(null)

const rules = {
  required: (v: string) => !!v || '此字段为必填项'
}

async function handleLogin() {
  const { valid } = await formRef.value.validate()
  if (!valid) return

  loading.value = true
  errorMessage.value = ''

  try {
    await userStore.login(username.value, password.value)
    router.push('/profile')
  } catch (error: any) {
    errorMessage.value = error.response?.data?.error || '登录失败，请检查用户名和密码'
  } finally {
    loading.value = false
  }
}
</script>
