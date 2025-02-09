<template>
  <div class="user-card">
    <span class="user-icon">👤</span>
    <span class="username">{{ user.username }}</span>
  </div>
</template>

<script>
export default {
  props: {
    user: Object,
  },
};
</script>

<style scoped>
.user-card {
  display: flex;
  align-items: center;
  justify-content: space-between;
  background: #d6c3ff;
  padding: 15px;
  border-radius: 10px;
  font-size: 18px;
  font-weight: bold;
  cursor: pointer;
}
</style>