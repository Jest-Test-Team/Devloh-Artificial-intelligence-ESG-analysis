import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:devlog_ai/screens/logs/logs_screen.dart';
import 'package:intl/intl.dart';

class LogFilterDialog extends ConsumerStatefulWidget {
  final List<String> availableServices;

  const LogFilterDialog({
    required this.availableServices,
    Key? key,
  }) : super(key: key);

  @override
  ConsumerState<LogFilterDialog> createState() => _LogFilterDialogState();
}

class _LogFilterDialogState extends ConsumerState<LogFilterDialog> {
  late Set<String> _selectedLevels;
  late Set<String> _selectedServices;
  late DateTime? _startDate;
  late DateTime? _endDate;
  
  final List<String> _allLevels = ['ERROR', 'WARN', 'INFO', 'DEBUG'];
  
  @override
  void initState() {
    super.initState();
    final filter = ref.read(logFilterProvider);
    _selectedLevels = Set.from(filter.selectedLevels);
    _selectedServices = Set.from(filter.selectedServices);
    _startDate = filter.startDate;
    _endDate = filter.endDate;
  }
  
  @override
  Widget build(BuildContext context) {
    return AlertDialog(
      title: const Text('過濾日誌'),
      content: SingleChildScrollView(
        child: Column(
          mainAxisSize: MainAxisSize.min,
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            // 日誌級別過濾器
            const Text(
              '日誌級別',
              style: TextStyle(
                fontWeight: FontWeight.bold,
              ),
            ),
            const SizedBox(height: 8),
            Wrap(
              spacing: 8,
              children: _allLevels.map((level) {
                final isSelected = _selectedLevels.contains(level);
                return FilterChip(
                  label: Text(level),
                  selected: isSelected,
                  onSelected: (selected) {
                    setState(() {
                      if (selected) {
                        _selectedLevels.add(level);
                      } else {
                        _selectedLevels.remove(level);
                      }
                    });
                  },
                );
              }).toList(),
            ),
            const SizedBox(height: 16),
            
            // 服務過濾器
            const Text(
              '服務',
              style: TextStyle(
                fontWeight: FontWeight.bold,
              ),
            ),
            const SizedBox(height: 8),
            Wrap(
              spacing: 8,
              children: widget.availableServices.map((service) {
                final isSelected = _selectedServices.contains(service);
                return FilterChip(
                  label: Text(service),
                  selected: isSelected,
                  onSelected: (selected) {
                    setState(() {
                      if (selected) {
                        _selectedServices.add(service);
                      } else {
                        _selectedServices.remove(service);
                      }
                    });
                  },
                );
              }).toList(),
            ),
            const SizedBox(height: 16),
            
            // 日期範圍過濾器
            const Text(
              '日期範圍',
              style: TextStyle(
                fontWeight: FontWeight.bold,
              ),
            ),
            const SizedBox(height: 8),
            Row(
              children: [
                Expanded(
                  child: InkWell(
                    onTap: () => _selectDate(true),
                    child: InputDecorator(
                      decoration: const InputDecoration(
                        labelText: '開始日期',
                        border: OutlineInputBorder(),
                      ),
                      child: Text(
                        _startDate != null
                            ? DateFormat('yyyy-MM-dd HH:mm').format(_startDate!)
                            : '選擇開始日期',
                      ),
                    ),
                  ),
                ),
                const SizedBox(width: 16),
                Expanded(
                  child: InkWell(
                    onTap: () => _selectDate(false),
                    child: InputDecorator(
                      decoration: const InputDecoration(
                        labelText: '結束日期',
                        border: OutlineInputBorder(),
                      ),
                      child: Text(
                        _endDate != null
                            ? DateFormat('yyyy-MM-dd HH:mm').format(_endDate!)
                            : '選擇結束日期',
                      ),
                    ),
                  ),
                ),
              ],
            ),
            if (_startDate != null || _endDate != null) ...[
              const SizedBox(height: 8),
              TextButton.icon(
                onPressed: () {
                  setState(() {
                    _startDate = null;
                    _endDate = null;
                  });
                },
                icon: const Icon(Icons.clear),
                label: const Text('清除日期'),
              ),
            ],
          ],
        ),
      ),
      actions: [
        TextButton(
          onPressed: () {
            Navigator.pop(context);
          },
          child: const Text('取消'),
        ),
        TextButton(
          onPressed: () {
            ref.read(logFilterProvider.notifier).updateSelectedLevels(_selectedLevels);
            ref.read(logFilterProvider.notifier).updateSelectedServices(_selectedServices);
            ref.read(logFilterProvider.notifier).updateDateRange(_startDate, _endDate);
            Navigator.pop(context);
          },
          child: const Text('套用'),
        ),
      ],
    );
  }
  
  // 選擇日期
  Future<void> _selectDate(bool isStartDate) async {
    final initialDate = isStartDate ? _startDate ?? DateTime.now() : _endDate ?? DateTime.now();
    
    final date = await showDatePicker(
      context: context,
      initialDate: initialDate,
      firstDate: DateTime(2000),
      lastDate: DateTime(2100),
    );
    
    if (date == null) return;
    
    // ignore: use_build_context_synchronously
    final time = await showTimePicker(
      context: context,
      initialTime: TimeOfDay.fromDateTime(initialDate),
    );
    
    if (time == null) return;
    
    final selectedDateTime = DateTime(
      date.year,
      date.month,
      date.day,
      time.hour,
      time.minute,
    );
    
    setState(() {
      if (isStartDate) {
        _startDate = selectedDateTime;
      } else {
        _endDate = selectedDateTime;
      }
    });
  }
}
