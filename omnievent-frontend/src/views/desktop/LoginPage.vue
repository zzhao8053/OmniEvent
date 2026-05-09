<template>
  <v-container class="fill-height" fluid>
    <v-row align="center" justify="center">
      <v-col cols="12" sm="8" md="4">
        <v-card class="elevation-12">
          <v-card-title class="text-center">
            <h2>登录 OmniEvent</h2>
          </v-card-title>
          <v-card-text>
            <v-form ref="formRef" @submit.prevent="show2faInput ? handleVerify() : handleLogin()">
              <v-text-field
                v-model="username"
                label="用户名或邮箱"
                :rules="[rules.required]"
                :disabled="show2faInput || loading"
                variant="outlined"
                prepend-inner-icon="mdi-account"
                autocomplete="username"
              />
              <v-text-field
                v-model="password"
                label="密码"
                :type="showPassword ? 'text' : 'password'"
                :rules="[rules.required]"
                :disabled="show2faInput || loading"
                variant="outlined"
                prepend-inner-icon="mdi-lock"
                :append-inner-icon="showPassword ? 'mdi-eye-off' : 'mdi-eye'"
                @click:append-inner="showPassword = !showPassword"
                autocomplete="current-password"
              />
              <v-expand-transition>
                <v-text-field
                  v-if="show2faInput"
                  v-model="passcode"
                  label="验证码"
                  type="text"
                  :disabled="loading"
                  variant="outlined"
                  prepend-inner-icon="mdi-two-factor-authentication"
                  :append-inner-icon="twoFAVerifyType === 'passcode' ? 'mdi-help-circle-outline' : 'mdi-key'"
                  @click:append-inner="toggle2FAVerifyType"
                  :placeholder="twoFAVerifyType === 'passcode' ? '输入6位验证码' : '输入恢复代码'"
                  autocomplete="one-time-code"
                />
              </v-expand-transition>
              <v-alert v-if="errorMessage" type="error" class="mb-4">
                {{ errorMessage }}
              </v-alert>
              <v-btn
                type="submit"
                color="primary"
                block
                :loading="loading"
                :disabled="show2faInput ? !passcode : !username || !password"
                size="large"
              >
                {{ show2faInput ? '验证' : '登录' }}
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
const passcode = ref('')
const showPassword = ref(false)
const loading = ref(false)
const errorMessage = ref('')
const formRef = ref<any>(null)
const tempToken = ref('')
const show2faInput = ref(false)
const twoFAVerifyType = ref<'passcode' | 'backupcode'>('passcode')

const rules = {
  required: (v: string) => !!v || '此字段为必填项'
}

function toggle2FAVerifyType() {
  twoFAVerifyType.value = twoFAVerifyType.value === 'passcode' ? 'backupcode' : 'passcode'
}

async function handleLogin() {
  const { valid } = await formRef.value.validate()
  if (!valid) return

  loading.value = true
  errorMessage.value = ''

  try {
    const response: any = await userStore.login(username.value, password.value)
    if (response.need2FA) {
      tempToken.value = response.token
      show2faInput.value = true
      loading.value = false
      return
    }
    router.push('/profile')
  } catch (error: any) {
    errorMessage.value = error.response?.data?.error || '登录失败，请检查用户名和密码'
  } finally {
    loading.value = false
  }
}

async function handleVerify() {
  if (!passcode.value) {
    errorMessage.value = '请输入验证码'
    return
  }

  loading.value = true
  errorMessage.value = ''

  try {
    await userStore.verify2FA(tempToken.value, passcode.value)
    router.push('/profile')
  } catch (error: any) {
    errorMessage.value = error.response?.data?.error || '验证失败'
  } finally {
    loading.value = false
  }
}
</script>
