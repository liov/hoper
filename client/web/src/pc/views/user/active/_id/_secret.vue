<template>
  <div />
</template>

<script>
import axios from 'axios'
export default {
  mounted() {
    const vm = this
    this.$axios
      .$get(
        `/api/user/active/${this.$route.params.id}/${this.$route.params.secret}`
      )
      .then((res) => {
        vm.$message.info(res.msg)
        if (res.code === 200) {
          vm.$router.push({ path: '/user/login?email=' + res.data })
        }
      })
  }
}
</script>

<style scoped></style>
