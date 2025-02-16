<template>
  <div v-if="message" class="error-box">
    <p>{{ message }}</p>
  </div>
</template>

<script>
export default {
  props: {
    message: {
      type: String,
      required: true
    }
  }
};
</script>

<style scoped>
.error-box {
  background-color: #ffdddd;
  color: #d8000c;
  padding: 10px;
  border: 1px solid #d8000c;
  border-radius: 5px;
  margin: 10px 0;
  text-align: center;
}
</style>
