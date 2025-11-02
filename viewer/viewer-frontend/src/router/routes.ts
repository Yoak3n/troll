import type { RouteRecordRaw } from "vue-router"

const routes: RouteRecordRaw[] = [
    {
        path: "/",
        name: "root",
        component: () => import("../views/Layout.vue"),
        children: [
            {
                path: "/home", name: "home", component: () => import("../views/Home.vue"), children: [
                    { path: "/topics", name: "topics", component: () => import("../views/Topics.vue") },
                    { path: "/videos", name: "videos", component: () => import("../views/Videos.vue") },
                    { path: "/user", name: "user", component: () => import("../views/User.vue") },
                ]
            },
            { path: "/console", name: "console", component: () => import("../views/Console.vue") },
            { path: "/setting", name: "setting", component: () => import("../views/Setting.vue") },
            { path: "/about", name: "about", component: () => import("../views/About.vue") },
        ]
    }
]
export default routes