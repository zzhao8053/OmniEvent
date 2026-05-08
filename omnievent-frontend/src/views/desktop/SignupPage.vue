<template>
  <v-container class="fill-height" fluid>
    <v-row align="center" justify="center">
      <v-col cols="12" sm="8" md="4">
        <v-card class="elevation-12">
          <v-card-title class="text-center">
            <h2>注册 OmniEvent</h2>
          </v-card-title>
          <v-card-text>
            <v-form ref="formRef" @submit.prevent="handleSignup">
              <v-text-field
                v-model="form.username"
                label="用户名"
                :rules="[rules.required, rules.usernameLength]"
                variant="outlined"
                prepend-inner-icon="mdi-account"
              />
              <v-text-field
                v-model="form.email"
                label="邮箱"
                :rules="[rules.required, rules.email]"
                variant="outlined"
                prepend-inner-icon="mdi-email"
              />
              <v-text-field
                v-model="form.nickname"
                label="昵称"
                :rules="[rules.required]"
                variant="outlined"
                prepend-inner-icon="mdi-card-account-details"
              />
              <v-text-field
                v-model="form.password"
                label="密码"
                :type="showPassword ? 'text' : 'password'"
                :rules="[rules.required, rules.passwordLength]"
                variant="outlined"
                prepend-inner-icon="mdi-lock"
                :append-inner-icon="showPassword ? 'mdi-eye-off' : 'mdi-eye'"
                @click:append-inner="showPassword = !showPassword"
              />
              <v-text-field
                v-model="confirmPassword"
                label="确认密码"
                :type="showPassword ? 'text' : 'password'"
                :rules="[rules.required, passwordMatch]"
                variant="outlined"
                prepend-inner-icon="mdi-lock-check"
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
                注册
              </v-btn>
            </v-form>
          </v-card-text>
          <v-card-actions>
            <v-spacer />
            <v-btn variant="text" @click="$router.push('/login')">
              已有账号？去登录
            </v-btn>
          </v-card-actions>
        </v-card>
      </v-col>
    </v-row>
  </v-container>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue'
import { useRouter } from 'vue-router'
import { useUserStore } from '@/stores/user'

const router = useRouter()
const userStore = useUserStore()

const form = reactive({
  username: '',
  email: '',
  nickname: '',
  password: ''
})
const confirmPassword = ref('')
const showPassword = ref(false)
const loading = ref(false)
const errorMessage = ref('')
const successMessage = ref('')
const formRef = ref<any>(null)

const rules = {
  required: (v: string) => !!v || '此字段为必填项',
  usernameLength: (v: string) => (v.length >= 3 && v.length <= 32) || '用户名长度应为3-32个字符',
  passwordLength: (v: string) => v.length >= 6 || '密码至少6个字符',
  email: (v: string) => /.+@.+\..+/.test(v) || '请输入有效的邮箱地址'
}

function passwordMatch(v: string) {
  return v === form.password || '两次输入的密码不一致'
}

async function handleSignup() {
  const { valid } = await formRef.value.validate()
  if (!valid) return

  loading.value = true
  errorMessage.value = ''
  successMessage.value = ''

  try {
    await userStore.register(form.username, form.email, form.nickname, form.password)
    successMessage.value = '注册成功！请登录'
    setTimeout(() => {
      router.push('/login')
    }, 1500)
  } catch (error: any) {
    errorMessage.value = error.response?.data?.error || '注册失败，请稍后重试'
  } finally {
    loading.value = false
  }
}
</script>
