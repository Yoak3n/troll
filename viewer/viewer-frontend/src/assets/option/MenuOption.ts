import { h, type Component } from 'vue';
import {
    ChatbubbleEllipsesOutline,
    CaretForwardCircleSharp,
    HomeOutline,
    PersonOutline,
    BowlingBallOutline,
    CogOutline,
    HeartOutline
} from '@vicons/ionicons5'
import { NIcon, type MenuOption } from 'naive-ui';
export const MenuOptions: MenuOption[] = [
    {
        label: '首页', key: 'home', icon: renderIcon(HomeOutline), children: [
            { label: '话题列表', key: 'topics', icon: renderIcon(ChatbubbleEllipsesOutline) },
            { label: '视频列表', key: 'videos', icon: renderIcon(CaretForwardCircleSharp) },
            { label: '用户查询', key: 'users', icon: renderIcon(PersonOutline) }
        ]
    },
    { label: '控制台', key: 'console', icon: renderIcon(BowlingBallOutline) },
    { label: '设置', key: 'setting', icon: renderIcon(CogOutline) },
    { label: '关于', key: "about", icon: renderIcon(HeartOutline) }
];

function renderIcon(icon: Component) {
    return () => h(NIcon, null, { default: () => h(icon) })
}