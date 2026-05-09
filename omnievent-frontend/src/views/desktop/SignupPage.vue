<template>
  <div class="layout-wrapper">
    <router-link to="/">
      <div class="auth-logo d-flex align-start gap-x-3">
        <img alt="logo" class="login-page-logo" :src="APPLICATION_LOGO_PATH" />
        <h1 class="font-weight-medium leading-normal text-2xl">OmniEvent</h1>
      </div>
    </router-link>
    <v-row no-gutters class="auth-wrapper">
      <v-col cols="12" md="4" class="auth-image-background d-none d-md-flex align-center justify-center position-relative">
        <div class="d-flex auth-img-footer" v-if="!isDarkMode">
          <v-img class="img-with-direction" src="img/desktop/background.svg"/>
        </div>
        <div class="d-flex auth-img-footer" v-if="isDarkMode">
          <v-img class="img-with-direction" src="img/desktop/background-dark.svg"/>
        </div>
        <div class="d-flex align-center justify-center w-100 pt-10">
          <v-img class="img-with-direction" max-width="320px" src="img/desktop/people2.svg" v-if="!isDarkMode"/>
          <v-img class="img-with-direction" max-width="320px" src="img/desktop/people2-dark.svg" v-else-if="isDarkMode"/>
        </div>
      </v-col>
      <v-col cols="12" md="8" class="auth-card d-flex align-center justify-center pa-10">
        <v-card variant="flat" class="mt-12 mt-sm-0 pt-sm-12 pt-md-0">
          <h4 class="text-h4 mb-1">注册</h4>
          <p class="text-sm mt-2 mb-5">
            <span>已有账号？</span>
            <router-link class="ms-1" to="/login">点击登录</router-link>
          </p>
          <v-form ref="formRef" @submit.prevent="handleSignup">
            <v-row>
              <v-col cols="12" md="6">
                <v-text-field
                  type="text"
                  autocomplete="username"
                  autocapitalize="none"
                  autocorrect="off"
                  spellcheck="false"
                  inputmode="email"
                  :disabled="loading"
                  label="用户名"
                  placeholder="请输入用户名"
                  v-model="form.username"
                />
              </v-col>
              <v-col cols="12" md="6">
                <v-text-field
                  type="text"
                  autocomplete="nickname"
                  :disabled="loading"
                  label="昵称"
                  placeholder="请输入昵称"
                  v-model="form.nickname"
                />
              </v-col>
            </v-row>
            <v-row>
              <v-col cols="12" md="12">
                <v-text-field
                  type="email"
                  autocomplete="email"
                  :disabled="loading"
                  label="邮箱"
                  placeholder="请输入邮箱地址"
                  v-model="form.email"
                />
              </v-col>
            </v-row>
            <v-row>
              <v-col cols="12" md="6">
                <v-text-field
                  autocomplete="new-password"
                  type="password"
                  :disabled="loading"
                  label="密码"
                  placeholder="密码至少6个字符"
                  v-model="form.password"
                />
              </v-col>
              <v-col cols="12" md="6">
                <v-text-field
                  autocomplete="new-password"
                  type="password"
                  :disabled="loading"
                  label="确认密码"
                  placeholder="请再次输入密码"
                  v-model="confirmPassword"
                />
              </v-col>
            </v-row>

            <v-alert v-if="errorMessage" type="error" class="mb-4">
              {{ errorMessage }}
            </v-alert>
            <v-alert v-if="successMessage" type="success" class="mb-4">
              {{ successMessage }}
            </v-alert>

            <div class="d-flex justify-sm-space-between gap-4 flex-wrap mt-5">
              <v-btn
                color="primary"
                :loading="loading"
                size="large"
                @click="handleSignup"
              >
                注册
              </v-btn>
            </div>
          </v-form>
        </v-card>
      </v-col>
    </v-row>

    <snack-bar ref="snackbar" @update:show="onSnackbarShowStateChanged" />
  </div>
</template>

<script setup lang="ts">
import SnackBar from '@/components/desktop/SnackBar.vue';

import { ref, reactive, useTemplateRef } from 'vue';
import { useRouter } from 'vue-router';
import { useTheme } from 'vuetify';

import { useUserStore } from '@/stores/user'

import { ThemeType } from '@/core/theme.ts';
import { APPLICATION_LOGO_PATH } from '@/consts/asset.ts';

import {
    mdiArrowLeft,
    mdiArrowRight,
    mdiCheck
} from '@mdi/js';

type SnackBarType = InstanceType<typeof SnackBar>;

const router = useRouter();
const theme = useTheme();

const userStore = useUserStore();

const snackbar = useTemplateRef<SnackBarType>('snackbar');

const form = reactive({
  username: '',
  email: '',
  nickname: '',
  password: ''
})
const confirmPassword = ref('')
const loading = ref(false)
const errorMessage = ref('')
const successMessage = ref('')
const formRef = ref<any>(null)

const isDarkMode = computed<boolean>(() => theme.global.name.value === ThemeType.Dark);

function passwordMatch(v: string) {
  return v === form.password || '两次输入的密码不一致'
}

function validateForm(): boolean {
  if (!form.username || form.username.length < 3 || form.username.length > 32) {
    errorMessage.value = '用户名长度应为3-32个字符'
    return false
  }
  if (!form.nickname) {
    errorMessage.value = '昵称为必填项'
    return false
  }
  if (!/.+@.+\..+/.test(form.email)) {
    errorMessage.value = '请输入有效的邮箱地址'
    return false
  }
  if (form.password.length < 6) {
    errorMessage.value = '密码至少6个字符'
    return false
  }
  if (form.password !== confirmPassword.value) {
    errorMessage.value = '两次输入的密码不一致'
    return false
  }
  return true
}

async function handleSignup() {
  errorMessage.value = ''
  successMessage.value = ''

  if (!validateForm()) {
    return
  }

  loading.value = true

  try {
    await userStore.register(form.username, form.email, form.nickname, form.password)
    successMessage.value = '注册成功！请登录'
    setTimeout(() => {
      router.push('/login')
    }, 1500)
  } catch (error: any) {
    errorMessage.value = error.message || '注册失败，请稍后重试'
  } finally {
    loading.value = false
  }
}

function onSnackbarShowStateChanged(newValue: boolean): void {
  if (!newValue && successMessage.value) {
    router.replace('/login')
  }
}
</script>

<style scoped>
.auth-logo {
  padding: 24px 40px;
}

.login-page-logo {
  width: 40px;
  height: 40px;
}

.auth-wrapper {
  height: calc(100vh - 88px);
}

.auth-image-background {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  min-height: 100%;
}

.auth-card {
  background: rgb(var(--v-theme-surface));
}

.auth-img-footer {
  position: absolute;
  bottom: 0;
  width: 100%;
}

.img-with-direction {
  width: 100%;
}

.gap-x-3 {
  gap: 12px;
}
</style>
