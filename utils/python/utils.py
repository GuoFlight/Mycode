# 从列表中移除某元素
def rm_from_list(original_list, element_to_remove):
    return [element for element in original_list if element != element_to_remove]
