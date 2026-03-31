import request from "@/utils/request";
export type NoticeStatus = "draft" | "published";
export type NoticePriority = "normal" | "high" | "urgent";
export interface Notice {
  id: number;
  title: string;
  content: string;
  status: NoticeStatus;
  priority: NoticePriority;
  is_pinned: boolean;
  created_at: string;
  created_by?: number;
}

export interface NoticeFormData {
  title: string;
  content: string;
  status: NoticeStatus;
  is_pinned: boolean;
  priority: NoticePriority;
}
export interface NoticeListData {
  notices: Notice[];
  total: number;
  page: number;
  page_size: number;
  priority: NoticePriority;
}
export const noticeApi = {
  list(params?: {
    keyword?: string;
    status?: string;
    is_pinned?: boolean;
    priority?: NoticePriority;
    page?: number;
    page_size?: number;
  }) {
    return request.get<{ code: number; message: string; data: NoticeListData }>(
      "/notices",
      { params },
    );
  },
  create(data: NoticeFormData) {
    return request.post<{ code: number; message: string }>("/notices", data);
  },
  get(id: number) {
    return request.get<{ code: number; message: string; data: Notice }>(
      `/notices/${id}`,
    );
  },
  update(id: number, data: NoticeFormData) {
    return request.put<{ code: number; message: string; data: Notice }>(
      `/notices/${id}`,
      data,
    );
  },
  remove(id: number) {
    return request.delete<{ code: number; message: string }>(`/notices/${id}`);
  },
};
