import request from "../../utils/request";

const API = {
  STATS_URL: "/statistics/"
} as const;

export interface DashboardStats {
  topics: number;
  videos: number;
  users: number;
  comments: number;
}

export const fetchDashboardStats = () =>
  request.get<any, DashboardStats>(API.STATS_URL);