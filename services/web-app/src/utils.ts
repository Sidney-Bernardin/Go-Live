import { ref } from "vue";

export const setSessionID = async (sid: string): Promise<void> => {
    sessionStorage.setItem("sid", sid);
    // await setTimeout(() => { }, 1000);
}
export const getSessionID = (): string | null => sessionStorage.getItem("sid");
export const deleteSessionID = (): void => sessionStorage.removeItem("sid");

export const unexpectedErr = (err: any): void => {
    console.error(`An unexpected problem has occurred.`, err);
    alert(`An unexpected problem has occurred.`);
};

export const loader = () => {
    const loading = ref(false);

    const wrapLoad = async (fn: Promise<void>) => {
        loading.value = true;
        await fn.finally(() => (loading.value = false));
    };

    return {
        loading,
        wrapLoad,
    };
};
