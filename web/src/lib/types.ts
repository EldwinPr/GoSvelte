export interface Column {
	key: string;
	label: string;
	align?: 'left' | 'center' | 'right';
	width?: string;
	sortable?: boolean;
}

export interface PaginatedResult<T> {
	items: T[];
	total_count: number;
	page: number;
	page_size: number;
}
