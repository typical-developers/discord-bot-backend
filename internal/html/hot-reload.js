/**
 * This is not the best hot reloading method, ideally websockets would be better.
 * However.. I am not adding websockets to this project.
 * 
 * In the future this should just be proxied and separated so then it can be proxied.
 */

const versionUrl = '/html-version.json';

async function getVersion() {
    try {
        const res = await fetch(versionUrl);
        const data = await res.json();

        return data.version;
    } catch (e) {
        return await new Promise((res) => setTimeout(res(getVersion()), 1000));
    }
}

async function hotReload() {
    const currentVersion = await getVersion(); // initally get the version
    console.log('[hot-reload] Initial version:', currentVersion);

    while (true) {
        // then recheck it every so often.
        const newVersion = await getVersion();
        if (newVersion !== currentVersion) {
            console.log('[hot-reload] Version changed:', newVersion);
            location.reload();
        }

        await new Promise((res) => setTimeout(res, 2_500));
    }
}

(async () => await hotReload())();