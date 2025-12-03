import type { RouteRecordRaw } from "vue-router"

const routes: RouteRecordRaw[] = [
    {
        path: "/",
        name: "root",
        component: () => import("../views/Layout.vue"),
        redirect: '/home',
        children: [
            {
                path: "home", name: "home", component: () => import("../views/Home.vue"), children: [
                    {
                        path: "topics", name: "topics", component: () => import("../views/Topics.vue"), children: [
                            {
                                path: "topic", name: "topic", component: () => import("../views/TopicDetail.vue"), children: [
                                    { path: "video", name: "video", component: () => import("../views/VideoDetail.vue")}
                                ]
                            }
                        ]
                    },
                    { path: "videos", name: "videos", component: () => import("../views/Videos.vue") },
                    {
                        path: "users", name: "users", component: () => import("../views/Users.vue"), redirect: '/home/users/user-search', children: [
                            { path: "user-search", name: "user-search", component: () => import('../views/UserSearch.vue') },
                            { path: "user", name: "user", component: () => import('../views/UserDetail.vue') }
                        ]
                    },
                ]
            },
            { path: "console", name: "console", component: () => import("../views/Console.vue") },
            { path: "setting", name: "setting", component: () => import("../views/Setting.vue") },
            { path: "about", name: "about", component: () => import("../views/About.vue") },
        ]
    }
]
export default routes