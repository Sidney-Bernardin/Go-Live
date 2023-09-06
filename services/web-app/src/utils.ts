import { ref } from "vue";

export const setSessionID = (sid: string) => sessionStorage.setItem("sid", sid);
export const getSessionID = () => sessionStorage.getItem("sid");
export const deleteSessionID = () => sessionStorage.removeItem("sid");

export const unexpectedErr = (err: any) => {
  console.error(`An unexpected problem has occurred.`, err);
  alert(`An unexpected problem has occurred.`);
};

export const problems = {
  unauthorized: "unauthorized",
  invalidSignupInfo: "invalid_signup_info",
};

export function loader() {
  const loading = ref(false);

  const wrapLoad = async (fn: Promise<void>) => {
    loading.value = true;
    await fn.finally(() => (loading.value = false));
  };

  return {
    loading,
    wrapLoad,
  };
}
